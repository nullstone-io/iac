package yaml

import "gopkg.in/nullstone-io/go-api-client.v0/types"

type SubdomainConfiguration struct {
	BlockConfiguration `yaml:",inline" json:",inline"`

	SubdomainName string `yaml:"subdomain_name,omitempty" json:"subdomainName"`

	// DnsName
	// Deprecated: Use SubdomainName instead
	DnsName string `yaml:"dns_name,omitempty" json:"dnsName"`
}

func SubdomainConfigurationFromWorkspaceConfig(stackId, envId int64, config types.WorkspaceConfig) SubdomainConfiguration {
	return SubdomainConfiguration{
		BlockConfiguration: BlockConfigurationFromWorkspaceConfig(stackId, envId, config),
		SubdomainName:      config.Extra.SubdomainName,
	}
}
