package protocol

import (
	"encoding/json"
	"os"
	"sort"
	"time"

	deviceruntime "github.com/openclaw/openclaw-node/internal/device/runtime"
)

var startTime = time.Now()

func init() {
	RegisterHandler("device.describe", handleDeviceDescribe)
	RegisterHandler("device.info", handleDeviceInfo)
	RegisterHandler("device.status", handleDeviceStatus)
	RegisterHandler("device.health", handleDeviceHealth)
	RegisterHandler("device.permissions", handleDevicePermissions)
}

func handleDeviceDescribe(params json.RawMessage) (*InvokeResult, error) {
	rt := deviceruntime.Default()
	commands := append([]string{"device.describe", "device.info", "device.status", "device.health", "device.permissions"}, rt.Registry().Commands()...)
	sort.Strings(commands)
	return NewInvokeResultOK(map[string]interface{}{
		"caps":        rt.Registry().CapabilityNames(),
		"commands":    commands,
		"permissions": rt.Registry().Permissions(),
	})
}

func handleDeviceInfo(params json.RawMessage) (*InvokeResult, error) {
	hostname, _ := os.Hostname()
	deviceID := ""
	if GlobalProtocol != nil && GlobalProtocol.Identity != nil {
		deviceID = GlobalProtocol.Identity.DeviceID
	}
	meta := deviceruntime.Default().Metadata()
	return NewInvokeResultOK(map[string]interface{}{
		"platform":        meta.Platform,
		"os":              meta.OS,
		"osVersion":       meta.OSVersion,
		"model":           meta.Model,
		"modelIdentifier": meta.ModelIdentifier,
		"version":         meta.Version,
		"deviceId":        deviceID,
		"hostname":        hostname,
	})
}

func handleDeviceStatus(params json.RawMessage) (*InvokeResult, error) {
	return NewInvokeResultOK(map[string]interface{}{
		"connected":        true,
		"gatewayConnected": true,
		"nodeConnected":    true,
		"uptimeMs":         time.Since(startTime).Milliseconds(),
		"capabilities":     deviceruntime.Default().Registry().Availability(),
	})
}

func handleDeviceHealth(params json.RawMessage) (*InvokeResult, error) {
	return NewInvokeResultOK(map[string]interface{}{
		"ok":     true,
		"checks": []map[string]interface{}{{"name": "identity", "ok": GlobalProtocol != nil && GlobalProtocol.Identity != nil}, {"name": "storage", "ok": true}, {"name": "network", "ok": true}},
		"errors": []interface{}{},
	})
}

func handleDevicePermissions(params json.RawMessage) (*InvokeResult, error) {
	return NewInvokeResultOK(deviceruntime.Default().Registry().Permissions())
}
