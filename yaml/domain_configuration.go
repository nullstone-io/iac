package yaml

type DomainConfiguration struct {
	ModuleSource        string            `yaml:"module" json:"module"`
	ModuleSourceVersion *string           `yaml:"module_version,omitempty" json:"moduleVersion"`
	DnsName             string            `yaml:"dns_name" json:"dnsName"`
	Variables           map[string]any    `yaml:"vars" json:"vars"`
	Connections         ConnectionTargets `yaml:"connections" json:"connections"`
	IsShared            bool              `yaml:"is_shared" json:"isShared"`
}
