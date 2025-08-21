package config

import (
	"github.com/nullstone-io/iac/core"
	"github.com/nullstone-io/iac/yaml"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

type SubdomainConfiguration struct {
	BlockConfiguration

	SubdomainName string `json:"subdomainName"`
}

func (s *SubdomainConfiguration) ToBlock(orgName string, stackId int64) types.Block {
	block := s.BlockConfiguration.ToBlock(orgName, stackId)
	block.DnsName = s.SubdomainName
	return block
}

func (s *SubdomainConfiguration) ApplyChangesTo(ic core.IacContext, updater core.WorkspaceConfigUpdater) error {
	if err := s.BlockConfiguration.ApplyChangesTo(ic, updater); err != nil {
		return err
	}
	updater.UpdateSubdomainName(s.SubdomainName)
	return nil
}

func convertSubdomainConfigurations(parsed map[string]yaml.SubdomainConfiguration) map[string]*SubdomainConfiguration {
	result := make(map[string]*SubdomainConfiguration)
	for name, value := range parsed {
		bc := blockConfigFromYaml(name, value.BlockConfiguration, BlockTypeSubdomain, types.CategorySubdomain)
		subdomainName := value.SubdomainName
		if subdomainName == "" {
			subdomainName = value.DnsName
		}
		result[name] = &SubdomainConfiguration{BlockConfiguration: *bc, SubdomainName: subdomainName}
	}
	return result
}
