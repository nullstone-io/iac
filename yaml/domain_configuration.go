package yaml

import "gopkg.in/nullstone-io/go-api-client.v0/types"

type DomainConfiguration struct {
	BlockConfiguration `yaml:",inline" json:",inline"`

	Dns DomainDnsConfiguration `yaml:"dns" json:"dns"`
}

type DomainDnsConfiguration struct {
	Template string `yaml:"template" json:"template"`
}

func DomainConfigurationFromWorkspaceConfig(stackId, envId int64, config types.WorkspaceConfig) DomainConfiguration {
	domainTemplate := ""
	if config.Extra.Domain != nil {
		domainTemplate = config.Extra.Domain.DomainNameTemplate
	}
	return DomainConfiguration{
		BlockConfiguration: BlockConfigurationFromWorkspaceConfig(stackId, envId, config),
		Dns: DomainDnsConfiguration{
			Template: domainTemplate,
		},
	}
}
