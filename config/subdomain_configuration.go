package config

// TODO: Implement DnsName in DesiredConfig
type SubdomainConfiguration struct {
	BlockConfiguration
	DnsName string `yaml:"dns_name,omitempty" json:"dnsName"`
}
