package yaml

// TODO: Implement DnsName in DesiredConfig
type SubdomainConfiguration struct {
	ModuleSource        string            `yaml:"module" json:"module"`
	ModuleSourceVersion *string           `yaml:"module_version,omitempty" json:"moduleVersion"`
	DnsName             string            `yaml:"dns_name,omitempty" json:"dnsName"`
	Variables           map[string]any    `yaml:"vars" json:"vars"`
	Connections         ConnectionTargets `yaml:"connections" json:"connections"`
}
