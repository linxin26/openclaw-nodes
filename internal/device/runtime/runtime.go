package runtime

import (
	goruntime "runtime"
	"sort"
	"sync"

	"github.com/openclaw/openclaw-node/internal/config"
	capcalendar "github.com/openclaw/openclaw-node/internal/device/capabilities/calendar"
	capcamera "github.com/openclaw/openclaw-node/internal/device/capabilities/camera"
	caplocation "github.com/openclaw/openclaw-node/internal/device/capabilities/location"
	capmotion "github.com/openclaw/openclaw-node/internal/device/capabilities/motion"
	capnotifications "github.com/openclaw/openclaw-node/internal/device/capabilities/notifications"
	capphotos "github.com/openclaw/openclaw-node/internal/device/capabilities/photos"
	capscreen "github.com/openclaw/openclaw-node/internal/device/capabilities/screen"
	capsms "github.com/openclaw/openclaw-node/internal/device/capabilities/sms"
	"github.com/openclaw/openclaw-node/internal/device/core"
)

type Metadata struct {
	Platform        string
	OS              string
	OSVersion       string
	Model           string
	ModelIdentifier string
	Version         string
}

type ProviderSet struct {
	Camera        capcamera.Provider
	Photos        capphotos.Provider
	Screen        capscreen.Provider
	Location      caplocation.Provider
	Notifications capnotifications.Provider
	Calendar      capcalendar.Provider
}

type Runtime struct {
	cfg           *config.Config
	platform      string
	metadata      Metadata
	registry      *core.Registry
	camera        *capcamera.Service
	photos        *capphotos.Service
	screen        *capscreen.Service
	location      *caplocation.Service
	notifications *capnotifications.Service
	calendar      *capcalendar.Service
	motion        *capmotion.Service
	sms           *capsms.Service
}

var (
	current *Runtime
	mu      sync.RWMutex
)

func Default() *Runtime {
	mu.RLock()
	instance := current
	mu.RUnlock()
	if instance != nil {
		return instance
	}
	return MustBootstrap(config.Default())
}

func MustBootstrap(cfg *config.Config) *Runtime {
	instance := New(cfg, goruntime.GOOS)
	mu.Lock()
	current = instance
	mu.Unlock()
	return instance
}

func New(cfg *config.Config, platform string) *Runtime {
	if cfg == nil {
		cfg = config.Default()
	}
	set := BuildProviderSet(platform, cfg)
	rt := &Runtime{
		cfg:           cfg,
		platform:      platform,
		metadata:      metadataForPlatform(platform),
		camera:        capcamera.NewService(set.Camera),
		photos:        capphotos.NewService(set.Photos),
		screen:        capscreen.NewService(set.Screen),
		location:      caplocation.NewService(set.Location),
		notifications: capnotifications.NewService(set.Notifications),
		calendar:      capcalendar.NewService(set.Calendar),
		motion:        capmotion.NewService(),
		sms:           capsms.NewService(),
	}
	rt.registry = core.NewRegistry(rt.collectStates())
	return rt
}

func (r *Runtime) Metadata() Metadata                       { return r.metadata }
func (r *Runtime) Registry() *core.Registry                 { return r.registry }
func (r *Runtime) Camera() *capcamera.Service               { return r.camera }
func (r *Runtime) Photos() *capphotos.Service               { return r.photos }
func (r *Runtime) Screen() *capscreen.Service               { return r.screen }
func (r *Runtime) Location() *caplocation.Service           { return r.location }
func (r *Runtime) Notifications() *capnotifications.Service { return r.notifications }
func (r *Runtime) Calendar() *capcalendar.Service           { return r.calendar }
func (r *Runtime) Motion() *capmotion.Service               { return r.motion }
func (r *Runtime) SMS() *capsms.Service                     { return r.sms }
func (r *Runtime) Platform() string                         { return r.platform }

func (r *Runtime) collectStates() []core.CapabilityState {
	states := []core.CapabilityState{
		r.camera.State(r.isEnabled("camera")),
		r.photos.State(r.isEnabled("photos")),
		r.screen.State(r.isEnabled("screen")),
		r.location.State(r.isEnabled("location")),
		r.notifications.State(r.isEnabled("notifications")),
		r.calendar.State(r.isEnabled("calendar")),
		r.motion.State(r.isEnabled("motion")),
		r.sms.State(r.isEnabled("sms")),
	}
	sort.Slice(states, func(i, j int) bool { return states[i].Descriptor.Name < states[j].Descriptor.Name })
	return states
}

func (r *Runtime) isEnabled(name string) bool {
	if r.cfg == nil || r.cfg.Capabilities == nil {
		return false
	}
	return r.cfg.Capabilities[name]
}

func metadataForPlatform(platform string) Metadata {
	osName := platform
	switch platform {
	case "windows":
		osName = "Windows"
	case "darwin":
		osName = "macOS"
	case "linux":
		osName = "Linux"
	}
	return Metadata{
		Platform:        platform,
		OS:              osName,
		OSVersion:       "",
		Model:           "OpenClaw Node",
		ModelIdentifier: "openclaw-node-" + platform,
		Version:         "0.1.0",
	}
}
