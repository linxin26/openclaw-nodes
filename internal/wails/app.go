package wails

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	goruntime "runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/openclaw/openclaw-node/internal/config"
	appcrypto "github.com/openclaw/openclaw-node/internal/crypto"
	deviceruntime "github.com/openclaw/openclaw-node/internal/device/runtime"
	"github.com/openclaw/openclaw-node/internal/protocol"
)

var appInstance *App

type App struct {
	ctx      context.Context
	client   *protocol.Client
	config   *config.Config
	cfgPath  string
	identity *appcrypto.Identity
	dataDir  string

	mu             sync.RWMutex
	status         Status
	connectedAt    int64
	retryCount     int
	retryDelayMs   int64
	activityLog    []*ActivityEntry
	logBuffer      []*LogEntry
	connectCancel  context.CancelFunc
	lastStatusText string
}

func NewApp(dataDir string, identity *appcrypto.Identity, cfg *config.Config, client *protocol.Client) *App {
	if appInstance != nil {
		return appInstance
	}

	appInstance = &App{
		client:      client,
		config:      cfg,
		cfgPath:     filepath.Join(dataDir, "config.yaml"),
		identity:    identity,
		dataDir:     dataDir,
		status:      StatusOffline,
		activityLog: make([]*ActivityEntry, 0, 32),
		logBuffer:   make([]*LogEntry, 0, 256),
	}
	if client != nil {
		client.SetServerURL(normalizeGateway(cfg.Gateway, cfg.Port))
		client.SetToken(strings.TrimSpace(cfg.Token))
	}
	return appInstance
}

func GetApp() *App {
	return appInstance
}

func (a *App) WailsInit(ctx context.Context) {
	a.ctx = ctx
	a.wireClientCallbacks()
	a.pushLog("info", "Wails app initialized")
}

func (a *App) GetStatus() *ConnectionStatus {
	a.mu.RLock()
	status := a.status
	connectedAt := a.connectedAt
	retryCount := a.retryCount
	retryDelayMs := a.retryDelayMs
	a.mu.RUnlock()

	if a.client != nil && a.client.IsConnected() {
		status = StatusConnected
	}

	uptimeMs := int64(0)
	if connectedAt > 0 {
		uptimeMs = time.Now().UnixMilli() - connectedAt
	}
	capabilities := a.runtimeCapabilityStates()

	return &ConnectionStatus{
		Status:       status,
		Gateway:      normalizeGateway(a.config.Gateway, a.config.Port),
		TLS:          a.config.TLS,
		UptimeMs:     uptimeMs,
		RetryCount:   retryCount,
		RetryDelayMs: retryDelayMs,
		ProtocolV:    protocol.ProtocolVersion,
		Capabilities: capabilities,
	}
}

func (a *App) GetDeviceInfo() *DeviceInfo {
	hostname, _ := os.Hostname()
	meta := deviceruntime.Default().Metadata()
	return &DeviceInfo{
		DeviceID: a.identity.DeviceID,
		Platform: meta.Platform,
		Hostname: hostname,
		Mode:     "node",
		Version:  meta.Version,
	}
}

func (a *App) GetConfig() *Config {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return &Config{
		Gateway:           a.config.Gateway,
		Port:              a.config.Port,
		Token:             a.config.Token,
		TLS:               a.config.TLS,
		Discovery:         a.config.Discovery,
		Capabilities:      cloneCapabilities(a.config.Capabilities),
		CapabilityOptions: cloneCapabilityOptions(a.config.CapabilityOptions),
	}
}

func (a *App) SaveConfig(next *Config) error {
	if next == nil {
		return fmt.Errorf("config is required")
	}

	a.mu.Lock()
	a.config.Gateway = strings.TrimSpace(next.Gateway)
	a.config.Port = next.Port
	a.config.Token = strings.TrimSpace(next.Token)
	a.config.TLS = next.TLS
	a.config.Discovery = next.Discovery
	a.config.Capabilities = cloneCapabilities(next.Capabilities)
	a.config.CapabilityOptions = toConfigCapabilityOptions(next.CapabilityOptions)
	if a.client != nil {
		a.client.SetServerURL(normalizeGateway(a.config.Gateway, a.config.Port))
		a.client.SetToken(a.config.Token)
	}
	err := a.config.Save(a.cfgPath)
	a.mu.Unlock()
	if err != nil {
		return err
	}
	deviceruntime.MustBootstrap(a.config)

	a.pushActivity("configuration saved", "info")
	a.pushLog("info", "Configuration updated from GUI")
	EmitConfigChange(a.ctx, a.GetConfig())
	EmitStatusChange(a.ctx, a.GetStatus())
	return nil
}

func (a *App) Connect() error {
	if a.client == nil {
		return fmt.Errorf("protocol client is not initialized")
	}
	if strings.TrimSpace(a.config.Gateway) == "" {
		return fmt.Errorf("gateway is required")
	}

	a.mu.Lock()
	if a.connectCancel != nil {
		a.mu.Unlock()
		return nil
	}
	ctx, cancel := context.WithCancel(context.Background())
	a.connectCancel = cancel
	a.status = StatusConnecting
	a.retryCount++
	a.lastStatusText = "connecting"
	a.mu.Unlock()

	a.pushActivity("connecting to gateway", "info")
	a.pushLog("info", "Connecting to gateway")
	EmitStatusChange(a.ctx, a.GetStatus())

	go func() {
		err := a.client.Connect(ctx)
		if err != nil && !isContextCancellation(err) {
			a.mu.Lock()
			a.status = StatusError
			a.connectedAt = 0
			a.connectCancel = nil
			a.lastStatusText = err.Error()
			a.mu.Unlock()
			a.pushActivity("connection failed: "+err.Error(), "error")
			a.pushLog("error", "Connection failed: "+err.Error())
			EmitStatusChange(a.ctx, a.GetStatus())
			return
		}

		a.mu.Lock()
		a.connectCancel = nil
		if !a.client.IsConnected() && a.status != StatusConnected {
			a.status = StatusOffline
		}
		a.mu.Unlock()
		EmitStatusChange(a.ctx, a.GetStatus())
	}()

	return nil
}

func (a *App) Disconnect() error {
	a.mu.Lock()
	if a.connectCancel != nil {
		a.connectCancel()
		a.connectCancel = nil
	}
	a.status = StatusOffline
	a.connectedAt = 0
	a.lastStatusText = "disconnected"
	a.mu.Unlock()

	if a.client != nil {
		a.client.Disconnect()
	}

	a.pushActivity("disconnected from gateway", "warn")
	a.pushLog("warn", "Disconnected from gateway")
	EmitStatusChange(a.ctx, a.GetStatus())
	return nil
}

func (a *App) TestConnection() (*TestResult, error) {
	start := time.Now()
	endpoint := websocketEndpoint(normalizeGateway(a.config.Gateway, a.config.Port), a.config.TLS)

	dialer := websocket.Dialer{HandshakeTimeout: 5 * time.Second}
	conn, _, err := dialer.Dial(endpoint, nil)
	if err != nil {
		return &TestResult{
			Success:   false,
			LatencyMs: time.Since(start).Milliseconds(),
			Error:     err.Error(),
		}, nil
	}
	_ = conn.Close()

	return &TestResult{
		Success:   true,
		LatencyMs: time.Since(start).Milliseconds(),
	}, nil
}

func (a *App) GetCapabilities() []*CapabilityInfo {
	return capabilityCatalog()
}

func (a *App) SetCapability(key string, enabled bool) error {
	a.mu.Lock()
	if a.config.Capabilities == nil {
		a.config.Capabilities = map[string]bool{}
	}
	a.config.Capabilities[key] = enabled
	err := a.config.Save(a.cfgPath)
	a.mu.Unlock()
	if err != nil {
		return err
	}
	deviceruntime.MustBootstrap(a.config)

	info := a.lookupCapability(key)
	info.Enabled = enabled
	state := a.runtimeCapabilityStates()[key]
	info.Available = state.Available
	info.Permission = state.Permission
	info.Reason = state.Reason
	a.pushActivity(fmt.Sprintf("capability %s set to %t", key, enabled), "info")
	EmitCapabilityChange(a.ctx, info)
	EmitConfigChange(a.ctx, a.GetConfig())
	EmitStatusChange(a.ctx, a.GetStatus())
	return nil
}

func (a *App) InvokeCommand(method string, params map[string]interface{}) (*InvokeResult, error) {
	start := time.Now()
	raw, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req := protocol.InvokeRequest{
		ID:      fmt.Sprintf("gui-%d", time.Now().UnixMilli()),
		NodeID:  a.identity.DeviceID,
		Command: method,
		Params:  raw,
	}
	result := protocol.Dispatch(req)
	if result == nil {
		return nil, fmt.Errorf("invoke returned nil result")
	}

	duration := time.Since(start).Milliseconds()
	if result.OK {
		payload := map[string]interface{}{}
		if len(result.Payload) > 0 {
			if err := json.Unmarshal(result.Payload, &payload); err != nil {
				payload["raw"] = string(result.Payload)
			}
		}
		a.pushActivity("invoke success: "+method, "info")
		a.pushLog("info", "Invoke succeeded: "+method)
		EmitInvokeComplete(a.ctx, method, true, duration)
		return &InvokeResult{
			Success:    true,
			Data:       payload,
			DurationMs: duration,
		}, nil
	}

	errMsg := "invoke failed"
	if result.Error != nil {
		errMsg = result.Error.Message
	}
	a.pushActivity("invoke failed: "+method, "error")
	a.pushLog("error", fmt.Sprintf("Invoke failed: %s (%s)", method, errMsg))
	EmitInvokeComplete(a.ctx, method, false, duration)
	return &InvokeResult{
		Success:    false,
		Error:      errMsg,
		DurationMs: duration,
	}, nil
}

func (a *App) GetLogs(filter *LogFilter) []*LogEntry {
	a.mu.RLock()
	entries := make([]*LogEntry, len(a.logBuffer))
	copy(entries, a.logBuffer)
	a.mu.RUnlock()
	return filterLogs(entries, filter)
}

func (a *App) GetRecentActivity() []*ActivityEntry {
	a.mu.RLock()
	defer a.mu.RUnlock()

	items := make([]*ActivityEntry, len(a.activityLog))
	copy(items, a.activityLog)
	return items
}

func (a *App) GetAbout() *AboutInfo {
	info, _ := debug.ReadBuildInfo()
	hostname, _ := os.Hostname()

	goVersion := goruntime.Version()
	if info != nil && info.GoVersion != "" {
		goVersion = info.GoVersion
	}

	return &AboutInfo{
		DeviceID:        a.identity.DeviceID,
		PublicKey:       a.identity.PublicKeyBase64(),
		Version:         "0.1.0",
		Platform:        goruntime.GOOS,
		Hostname:        hostname,
		GoVersion:       goVersion,
		Arch:            goruntime.GOARCH,
		DataDir:         a.dataDir,
		ProtocolVersion: protocol.ProtocolVersion,
	}
}

func (a *App) OpenPath(path string) error {
	return exec.Command("explorer", path).Start()
}

func (a *App) SaveFileToDisk(base64Data, filename string) error {
	if filename == "" {
		return fmt.Errorf("filename is required")
	}
	payload := base64Data
	if idx := strings.Index(payload, ","); strings.Contains(payload, "base64,") && idx >= 0 {
		payload = payload[idx+1:]
	}
	data, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		data, err = base64.RawStdEncoding.DecodeString(payload)
		if err != nil {
			return err
		}
	}

	target := filepath.Join(a.dataDir, filename)
	if err := os.WriteFile(target, data, 0o644); err != nil {
		return err
	}
	a.pushLog("info", "Saved file to "+target)
	return nil
}

func (a *App) wireClientCallbacks() {
	if a.client == nil {
		return
	}

	a.client.OnConnected = func(_ protocol.ConnectResponse) {
		a.mu.Lock()
		a.status = StatusConnected
		a.connectedAt = time.Now().UnixMilli()
		a.lastStatusText = "connected"
		a.mu.Unlock()

		a.pushActivity("connected to gateway", "info")
		a.pushLog("info", "Connected to gateway")
		EmitStatusChange(a.ctx, a.GetStatus())
	}

	a.client.OnDisconnected = func(reason string) {
		a.mu.Lock()
		a.status = StatusOffline
		a.connectedAt = 0
		a.lastStatusText = reason
		a.mu.Unlock()

		a.pushActivity("disconnected: "+reason, "warn")
		a.pushLog("warn", "Disconnected: "+reason)
		EmitStatusChange(a.ctx, a.GetStatus())
	}
}

func (a *App) pushActivity(event, level string) {
	entry := &ActivityEntry{
		Timestamp: time.Now().UnixMilli(),
		Event:     event,
		Level:     level,
	}

	a.mu.Lock()
	a.activityLog = append([]*ActivityEntry{entry}, a.activityLog...)
	if len(a.activityLog) > 64 {
		a.activityLog = a.activityLog[:64]
	}
	a.mu.Unlock()

	EmitActivity(a.ctx, entry)
}

func (a *App) pushLog(level, message string) {
	entry := &LogEntry{
		Timestamp: time.Now().UnixMilli(),
		Level:     level,
		Message:   message,
	}

	a.mu.Lock()
	a.logBuffer = append(a.logBuffer, entry)
	if len(a.logBuffer) > 5000 {
		a.logBuffer = a.logBuffer[len(a.logBuffer)-5000:]
	}
	a.mu.Unlock()

	EmitLog(a.ctx, entry)
}

func (a *App) lookupCapability(key string) *CapabilityInfo {
	for _, item := range capabilityCatalog() {
		if item.Key == key {
			return item
		}
	}
	state := a.runtimeCapabilityStates()[key]
	return &CapabilityInfo{Key: key, Name: key, Enabled: state.Enabled, Available: state.Available, Permission: state.Permission, Reason: state.Reason, Healthy: state.Available}
}

func capabilityCatalog() []*CapabilityInfo {
	rt := deviceruntime.Default()
	states := rt.Registry().States()
	result := make([]*CapabilityInfo, 0, len(states))
	for _, key := range rt.Registry().CapabilityNames() {
		state := states[key]
		result = append(result, &CapabilityInfo{
			Key:         key,
			Name:        state.Descriptor.DisplayName,
			Description: state.Descriptor.Description,
			Enabled:     state.Availability.Enabled,
			Available:   state.Availability.Available,
			Permission:  string(state.Permission),
			Reason:      state.Availability.Reason,
			Commands:    state.Descriptor.Commands,
			Healthy:     state.Availability.Available,
			Tier:        state.Descriptor.Tier,
		})
	}
	return result
}

func cloneCapabilities(input map[string]bool) map[string]bool {
	if input == nil {
		return map[string]bool{}
	}
	out := make(map[string]bool, len(input))
	for key, value := range input {
		out[key] = value
	}
	return out
}

func cloneCapabilityOptions(input map[string]config.CapabilityOption) map[string]CapabilityOption {
	if input == nil {
		return map[string]CapabilityOption{}
	}
	out := make(map[string]CapabilityOption, len(input))
	for key, value := range input {
		out[key] = CapabilityOption{Provider: value.Provider, Path: value.Path}
	}
	return out
}

func toConfigCapabilityOptions(input map[string]CapabilityOption) map[string]config.CapabilityOption {
	if input == nil {
		return map[string]config.CapabilityOption{}
	}
	out := make(map[string]config.CapabilityOption, len(input))
	for key, value := range input {
		out[key] = config.CapabilityOption{Provider: value.Provider, Path: value.Path}
	}
	return out
}

func (a *App) runtimeCapabilityStates() map[string]CapabilityState {
	rt := deviceruntime.Default()
	states := rt.Registry().States()
	out := make(map[string]CapabilityState, len(states))
	for key, state := range states {
		out[key] = CapabilityState{Enabled: state.Availability.Enabled, Available: state.Availability.Available, Permission: string(state.Permission), Reason: state.Availability.Reason}
	}
	return out
}

func filterLogs(entries []*LogEntry, filter *LogFilter) []*LogEntry {
	if filter == nil {
		return entries
	}

	levelSet := map[string]struct{}{}
	for _, level := range filter.Levels {
		levelSet[strings.ToLower(level)] = struct{}{}
	}

	search := strings.ToLower(strings.TrimSpace(filter.Search))
	filtered := make([]*LogEntry, 0, len(entries))
	for _, entry := range entries {
		if len(levelSet) > 0 {
			if _, ok := levelSet[strings.ToLower(entry.Level)]; !ok {
				continue
			}
		}
		if search != "" && !strings.Contains(strings.ToLower(entry.Message), search) {
			continue
		}
		filtered = append(filtered, entry)
	}

	start := filter.Offset
	if start < 0 {
		start = 0
	}
	if start >= len(filtered) {
		return []*LogEntry{}
	}

	end := len(filtered)
	if filter.Limit > 0 && start+filter.Limit < end {
		end = start + filter.Limit
	}
	return filtered[start:end]
}

func websocketEndpoint(gateway string, tls bool) string {
	gateway = strings.TrimSpace(gateway)
	if gateway == "" {
		return "ws://localhost:18789"
	}
	if strings.HasPrefix(gateway, "ws://") || strings.HasPrefix(gateway, "wss://") {
		return gateway
	}
	if strings.HasPrefix(gateway, "https://") {
		return "wss://" + strings.TrimPrefix(gateway, "https://")
	}
	if strings.HasPrefix(gateway, "http://") {
		return "ws://" + strings.TrimPrefix(gateway, "http://")
	}
	if tls {
		return "wss://" + gateway
	}
	return "ws://" + gateway
}

func NormalizeGatewayAddress(gateway string, port int) string {
	return normalizeGateway(gateway, port)
}

func normalizeGateway(gateway string, port int) string {
	gateway = strings.TrimSpace(gateway)
	if gateway == "" {
		return ""
	}

	if strings.Contains(gateway, "://") {
		parsed, err := url.Parse(gateway)
		if err != nil {
			return gateway
		}
		if parsed.Port() != "" || port <= 0 {
			return parsed.String()
		}
		parsed.Host = net.JoinHostPort(parsed.Hostname(), strconv.Itoa(port))
		return parsed.String()
	}

	if _, _, err := net.SplitHostPort(gateway); err == nil {
		return gateway
	}

	if port > 0 {
		return net.JoinHostPort(gateway, strconv.Itoa(port))
	}
	return gateway
}

func isContextCancellation(err error) bool {
	return err == context.Canceled || err == context.DeadlineExceeded
}
