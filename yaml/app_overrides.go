package yaml

type CapabilityOverrides []CapabilityOverride

type CapabilityOverride struct {
	ModuleSource        string            `yaml:"module,omitempty" json:"module"`
	ModuleSourceVersion *string           `yaml:"module_version,omitempty" json:"moduleVersion"`
	Variables           map[string]any    `yaml:"vars" json:"vars"`
	Connections         ConnectionTargets `yaml:"connections,omitempty" json:"connections"`
	Namespace           *string           `yaml:"namespace,omitempty" json:"namespace"`
}

type AppOverrides struct {
	Variables map[string]any `yaml:"vars"`

	EnvVariables map[string]string   `yaml:"environment"`
	Capabilities CapabilityOverrides `yaml:"capabilities"`
}
