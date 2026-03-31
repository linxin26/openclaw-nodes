package discovery

import (
	"context"
	"fmt"
	"net"
	"time"
)

type MDNS struct {
	ServiceName string
	Port        int
	Hostname    string
	TXT         map[string]string
}

func NewMDNS(deviceID string, port int) *MDNS {
	return &MDNS{
		ServiceName: fmt.Sprintf("openclaw-node-%s", deviceID[:8]),
		Port:        port,
		Hostname:    fmt.Sprintf("%s.local.", deviceID[:8]),
		TXT: map[string]string{
			"platform": "windows",
			"version":  "0.1.0",
		},
	}
}

func (m *MDNS) Register(ctx context.Context) error {
	// Simplified: just listen on the port
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", m.Port))
	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		ln.Close()
	}()

	return nil
}

func (m *MDNS) Discover(ctx context.Context, serviceType string) ([]net.Addr, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	addr, err := net.ResolveUDPAddr("udp", "224.0.0.251:5353")
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return nil, nil
}
