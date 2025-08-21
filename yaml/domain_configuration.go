package yaml

import "gopkg.in/nullstone-io/go-api-client.v0/types"

type DomainConfiguration struct {
	BlockConfiguration `yaml:",inline" json:",inline"`

	DomainName string `yaml:"domain_name,omitempty" json:"domainName"`
}

func DomainConfigurationFromWorkspaceConfig(stackId, envId int64, config types.WorkspaceConfig) DomainConfiguration {
	return DomainConfiguration{
		BlockConfiguration: BlockConfigurationFromWorkspaceConfig(stackId, envId, config),
		DomainName:         config.Extra.DomainName,
	}
}
