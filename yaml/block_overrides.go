package yaml

type BlockOverrides struct {
	ModuleSource        string         `yaml:"module" json:"module"`
	ModuleSourceVersion *string        `yaml:"module_version,omitempty" json:"moduleVersion,omitempty"`
	Variables           map[string]any `yaml:"vars"`
}
