package config

import (
	"github.com/nullstone-io/iac/core"
	"github.com/nullstone-io/iac/yaml"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

type DomainConfiguration struct {
	BlockConfiguration

	DomainNameTemplate *string `json:"domainNameTemplate"`
}

func (d *DomainConfiguration) ToBlock(orgName string, stackId int64) types.Block {
	block := d.BlockConfiguration.ToBlock(orgName, stackId)
	if d.DomainNameTemplate != nil {
		block.DnsName = *d.DomainNameTemplate
	}
	return block
}

func (d *DomainConfiguration) ApplyChangesTo(ic core.IacContext, updater core.WorkspaceConfigUpdater) error {
	if err := d.BlockConfiguration.ApplyChangesTo(ic, updater); err != nil {
		return err
	}
	updater.UpdateDomainName(d.DomainNameTemplate)
	return nil
}

func convertDomainConfigurations(parsed map[string]yaml.DomainConfiguration) map[string]*DomainConfiguration {
	result := make(map[string]*DomainConfiguration)
	for name, value := range parsed {
		bc := blockConfigFromYaml(name, value.BlockConfiguration, BlockTypeDomain, types.CategoryDomain)
		var dnsTemplate *string
		if value.Dns.Template != "" {
			dnsTemplate = &value.Dns.Template
		}
		result[name] = &DomainConfiguration{BlockConfiguration: *bc, DomainNameTemplate: dnsTemplate}
	}
	return result
}
