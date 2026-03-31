package protocol

import (
	"fmt"
	"strings"
)

func BuildAuthPayload(deviceId, clientId, clientMode, role string, scopes []string, signedAtMs int64, token, nonce, platform, deviceFamily string) string {
	scopeString := strings.Join(scopes, ",")
	platformNorm := normalizeField(platform)
	deviceFamilyNorm := normalizeField(deviceFamily)

	payload := fmt.Sprintf("v3|%s|%s|%s|%s|%s|%d|%s|%s|%s|%s",
		deviceId, clientId, clientMode, role, scopeString, signedAtMs, token, nonce, platformNorm, deviceFamilyNorm)

	return payload
}

func normalizeField(v string) string {
	if v == "" {
		return ""
	}
	var b strings.Builder
	for _, c := range strings.TrimSpace(v) {
		if c >= 'A' && c <= 'Z' {
			b.WriteRune(c + 32)
		} else {
			b.WriteRune(c)
		}
	}
	return b.String()
}

const ProtocolVersion = 3
