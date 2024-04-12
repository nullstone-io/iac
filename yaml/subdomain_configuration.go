package yaml

// TODO: Implement DnsName in DesiredConfig
type SubdomainConfiguration struct {
	BlockConfiguration `yaml:",inline" json:",inline"`

	DnsName string `yaml:"dns_name,omitempty" json:"dnsName"`
}
