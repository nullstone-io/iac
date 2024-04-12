package yaml

type CapabilityConfigurations []CapabilityConfiguration

type CapabilityConfiguration struct {
	ModuleSource        string            `yaml:"module,omitempty" json:"module"`
	ModuleSourceVersion *string           `yaml:"module_version,omitempty" json:"moduleVersion"`
	Variables           map[string]any    `yaml:"vars" json:"vars"`
	Connections         ConnectionTargets `yaml:"connections,omitempty" json:"connections"`
	Namespace           *string           `yaml:"namespace,omitempty" json:"namespace"`
}

type AppConfiguration struct {
	BlockConfiguration `yaml:",inline" json:",inline"`

	EnvVariables map[string]string        `yaml:"environment" json:"envVars"`
	Capabilities CapabilityConfigurations `yaml:"capabilities" json:"capabilities"`
}
