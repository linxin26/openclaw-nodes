package main

import (
	"sort"

	"github.com/openclaw/openclaw-node/internal/device/core"
	deviceruntime "github.com/openclaw/openclaw-node/internal/device/runtime"
)

type CommandRegistry struct {
	caps        []string
	cmds        []string
	permissions map[string]bool
}

func NewRegistry() *CommandRegistry {
	rt := deviceruntime.Default()
	permissions := map[string]bool{}
	for name, state := range rt.Registry().Permissions() {
		permissions[name] = state == core.PermissionGranted
	}
	cmds := append([]string{"device.describe", "device.info", "device.status", "device.health", "device.permissions"}, rt.Registry().Commands()...)
	sort.Strings(cmds)
	return &CommandRegistry{caps: rt.Registry().CapabilityNames(), cmds: cmds, permissions: permissions}
}

func (r *CommandRegistry) AllCommands() []string { return r.cmds }
func (r *CommandRegistry) AllCaps() []string     { return r.caps }
