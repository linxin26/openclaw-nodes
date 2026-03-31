package tray

import (
	deviceruntime "github.com/openclaw/openclaw-node/internal/device/runtime"
)

type Menu struct {
	capabilities []Capability
}

type Capability struct {
	Key         string
	Name        string
	Description string
	Enabled     bool
	Available   bool
	Permission  string
	Reason      string
}

func NewMenu() *Menu {
	return NewMenuWithRuntime(deviceruntime.Default())
}

func NewMenuWithRuntime(rt *deviceruntime.Runtime) *Menu {
	states := rt.Registry().States()
	items := make([]Capability, 0, len(states))
	for _, name := range rt.Registry().CapabilityNames() {
		state := states[name]
		items = append(items, Capability{Key: name, Name: state.Descriptor.DisplayName, Description: state.Descriptor.Description, Enabled: state.Availability.Enabled, Available: state.Availability.Available, Permission: string(state.Permission), Reason: state.Availability.Reason})
	}
	return &Menu{capabilities: items}
}

func (m *Menu) Capabilities() []Capability { return m.capabilities }

func (m *Menu) SetEnabled(name string, enabled bool) {
	for i := range m.capabilities {
		if m.capabilities[i].Key == name || m.capabilities[i].Name == name {
			m.capabilities[i].Enabled = enabled
			break
		}
	}
}
