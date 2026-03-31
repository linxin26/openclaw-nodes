package device

import (
	"testing"

	"github.com/openclaw/openclaw-node/internal/config"
	"github.com/openclaw/openclaw-node/internal/protocol"
)

func TestCameraHandlersAreRegisteredExplicitly(t *testing.T) {
	Bootstrap(config.Default())
	RegisterProtocolHandlers(protocol.RegisterHandler)

	for _, command := range []string{"camera.list", "camera.snap", "camera.clip"} {
		if _, ok := protocol.GetHandler(command); !ok {
			t.Fatalf("handler %q is not registered", command)
		}
	}
}
