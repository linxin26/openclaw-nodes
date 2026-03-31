package tray

import (
	"log"

	deviceruntime "github.com/openclaw/openclaw-node/internal/device/runtime"
)

type Tray struct {
	menu   *Menu
	icon   []byte
	status Status
}

type Status int

const (
	StatusOffline Status = iota
	StatusConnecting
	StatusConnected
	StatusError
)

func New() *Tray {
	return NewWithRuntime(deviceruntime.Default())
}

func NewWithRuntime(rt *deviceruntime.Runtime) *Tray {
	return &Tray{status: StatusOffline, menu: NewMenuWithRuntime(rt)}
}

func (t *Tray) Run() {
	log.Println("OpenClaw Node started")
	log.Println("Status: Offline")
	log.Println("Capabilities:")
	for _, cap := range t.menu.Capabilities() {
		log.Printf("  - %s: enabled=%v available=%v permission=%s reason=%s", cap.Name, cap.Enabled, cap.Available, cap.Permission, cap.Reason)
	}
	log.Println("Use --help for options")
}

func (t *Tray) SetStatus(s Status) {
	t.status = s
	title := [...]string{"Offline", "Connecting...", "Connected", "Error"}[s]
	log.Printf("Status changed to: %s", title)
}
