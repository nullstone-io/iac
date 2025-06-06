package yaml

import "gopkg.in/nullstone-io/go-api-client.v0/types"

// TODO: Implement DnsName in types.WorkspaceConfig
type DomainConfiguration struct {
	BlockConfiguration `yaml:",inline" json:",inline"`

	DnsName string `yaml:"dns_name,omitempty" json:"dnsName"`
}

func DomainConfigurationFromWorkspaceConfig(stackId, envId int64, block types.Block, config types.WorkspaceConfig) DomainConfiguration {
	return DomainConfiguration{
		BlockConfiguration: BlockConfigurationFromWorkspaceConfig(stackId, envId, config),
		DnsName:            block.DnsName,
	}
}
