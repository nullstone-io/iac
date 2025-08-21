package config

import (
	"github.com/nullstone-io/iac/core"
	"github.com/nullstone-io/iac/yaml"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

type DomainConfiguration struct {
	BlockConfiguration

	DomainName string `json:"domainName"`
}

func (d *DomainConfiguration) ToBlock(orgName string, stackId int64) types.Block {
	block := d.BlockConfiguration.ToBlock(orgName, stackId)
	return block
}

func (d *DomainConfiguration) ApplyChangesTo(ic core.IacContext, updater core.WorkspaceConfigUpdater) error {
	if err := d.BlockConfiguration.ApplyChangesTo(ic, updater); err != nil {
		return err
	}
	updater.UpdateDomainName(d.DomainName)
	return nil
}

func convertDomainConfigurations(parsed map[string]yaml.DomainConfiguration) map[string]*DomainConfiguration {
	result := make(map[string]*DomainConfiguration)
	for name, value := range parsed {
		bc := blockConfigFromYaml(name, value.BlockConfiguration, BlockTypeDomain, types.CategoryDomain)
		result[name] = &DomainConfiguration{BlockConfiguration: *bc, DomainName: value.DomainName}
	}
	return result
}
