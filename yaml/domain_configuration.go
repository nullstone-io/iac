package yaml

type DomainConfiguration struct {
	BlockConfiguration `yaml:",inline" json:",inline"`

	DnsName string `yaml:"dns_name" json:"dnsName"`
}
