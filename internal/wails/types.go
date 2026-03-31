package wails

type Status string

const (
	StatusOffline    Status = "offline"
	StatusConnecting Status = "connecting"
	StatusConnected  Status = "connected"
	StatusError      Status = "error"
)

type CapabilityState struct {
	Enabled    bool   `json:"enabled"`
	Available  bool   `json:"available"`
	Permission string `json:"permission"`
	Reason     string `json:"reason,omitempty"`
}

type CapabilityOption struct {
	Provider string `json:"provider,omitempty"`
	Path     string `json:"path,omitempty"`
}

type ConnectionStatus struct {
	Status       Status                     `json:"status"`
	Gateway      string                     `json:"gateway"`
	TLS          bool                       `json:"tls"`
	UptimeMs     int64                      `json:"uptimeMs"`
	RetryCount   int                        `json:"retryCount"`
	RetryDelayMs int64                      `json:"retryDelayMs"`
	ProtocolV    int                        `json:"protocolVersion"`
	Capabilities map[string]CapabilityState `json:"capabilities"`
}

type DeviceInfo struct {
	DeviceID string `json:"deviceId"`
	Platform string `json:"platform"`
	Hostname string `json:"hostname"`
	Mode     string `json:"mode"`
	Version  string `json:"version"`
}

type Config struct {
	Gateway           string                      `json:"gateway"`
	Port              int                         `json:"port"`
	Token             string                      `json:"token"`
	TLS               bool                        `json:"tls"`
	Discovery         string                      `json:"discovery"`
	Capabilities      map[string]bool             `json:"capabilities"`
	CapabilityOptions map[string]CapabilityOption `json:"capabilityOptions"`
}

type TestResult struct {
	Success   bool   `json:"success"`
	LatencyMs int64  `json:"latencyMs"`
	Error     string `json:"error"`
}

type CapabilityInfo struct {
	Key          string   `json:"key"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Enabled      bool     `json:"enabled"`
	Available    bool     `json:"available"`
	Permission   string   `json:"permission"`
	Reason       string   `json:"reason,omitempty"`
	Commands     []string `json:"commands"`
	Dependencies []string `json:"dependencies"`
	Healthy      bool     `json:"healthy"`
	Tier         int      `json:"tier"`
}

type InvokeResult struct {
	Success    bool                   `json:"success"`
	Data       map[string]interface{} `json:"data"`
	Error      string                 `json:"error"`
	DurationMs int64                  `json:"durationMs"`
}

type AboutInfo struct {
	DeviceID        string `json:"deviceId"`
	PublicKey       string `json:"publicKey"`
	Version         string `json:"version"`
	Platform        string `json:"platform"`
	Hostname        string `json:"hostname"`
	GoVersion       string `json:"goVersion"`
	Arch            string `json:"arch"`
	DataDir         string `json:"dataDir"`
	ProtocolVersion int    `json:"protocolVersion"`
}

type LogEntry struct {
	Timestamp int64  `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

type ActivityEntry struct {
	Timestamp int64  `json:"timestamp"`
	Event     string `json:"event"`
	Level     string `json:"level"`
}

type LogFilter struct {
	Levels []string `json:"levels"`
	Search string   `json:"search"`
	Limit  int      `json:"limit"`
	Offset int      `json:"offset"`
}
