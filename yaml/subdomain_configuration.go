package yaml

import "gopkg.in/nullstone-io/go-api-client.v0/types"

type SubdomainConfiguration struct {
	BlockConfiguration `yaml:",inline" json:",inline"`

	Dns SubdomainDnsConfiguration `yaml:"dns" json:"dns"`

	// DnsName
	// Deprecated: Use SubdomainName instead
	DnsName string `yaml:"dns_name,omitempty" json:"dnsName"`
}

type SubdomainDnsConfiguration struct {
	Template string `yaml:"template" json:"template"`
}

func SubdomainConfigurationFromWorkspaceConfig(stackId, envId int64, config types.WorkspaceConfig) SubdomainConfiguration {
	subdomainTemplate := ""
	if config.Extra.Subdomain != nil {
		subdomainTemplate = config.Extra.Subdomain.SubdomainNameTemplate
	}
	return SubdomainConfiguration{
		BlockConfiguration: BlockConfigurationFromWorkspaceConfig(stackId, envId, config),
		Dns: SubdomainDnsConfiguration{
			Template: subdomainTemplate,
		},
	}
}
