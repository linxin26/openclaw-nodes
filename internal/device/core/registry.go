package core

import "sort"

type Registry struct {
	states map[string]CapabilityState
}

func NewRegistry(states []CapabilityState) *Registry {
	items := make(map[string]CapabilityState, len(states))
	for _, state := range states {
		if state.Descriptor.Name == "" {
			continue
		}
		items[state.Descriptor.Name] = state
	}
	return &Registry{states: items}
}

func (r *Registry) CapabilityNames() []string {
	names := make([]string, 0, len(r.states))
	for name := range r.states {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func (r *Registry) Commands() []string {
	set := map[string]struct{}{}
	for _, state := range r.states {
		for _, command := range state.Descriptor.Commands {
			set[command] = struct{}{}
		}
	}
	commands := make([]string, 0, len(set))
	for command := range set {
		commands = append(commands, command)
	}
	sort.Strings(commands)
	return commands
}

func (r *Registry) Permissions() map[string]PermissionState {
	out := make(map[string]PermissionState, len(r.states))
	for name, state := range r.states {
		out[name] = state.Permission
	}
	return out
}

func (r *Registry) Availability() map[string]Availability {
	out := make(map[string]Availability, len(r.states))
	for name, state := range r.states {
		out[name] = state.Availability
	}
	return out
}

func (r *Registry) States() map[string]CapabilityState {
	out := make(map[string]CapabilityState, len(r.states))
	for name, state := range r.states {
		out[name] = state
	}
	return out
}
