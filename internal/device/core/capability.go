package core

type CapabilityDescriptor struct {
	Name        string
	Commands    []string
	Tier        int
	DisplayName string
	Description string
}

type Availability struct {
	Enabled   bool   `json:"enabled"`
	Available bool   `json:"available"`
	Reason    string `json:"reason,omitempty"`
}

type CapabilityState struct {
	Descriptor   CapabilityDescriptor
	Permission   PermissionState
	Availability Availability
}
