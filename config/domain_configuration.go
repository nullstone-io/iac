package config

import (
	"github.com/nullstone-io/iac/yaml"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

type DomainConfiguration struct {
	BlockConfiguration

	DnsName string `json:"dnsName"`
}

func (d *DomainConfiguration) ToBlock(orgName string, stackId int64) types.Block {
	block := d.BlockConfiguration.ToBlock(orgName, stackId)
	block.DnsName = d.DnsName
	return block
}

func convertDomainConfigurations(parsed map[string]yaml.DomainConfiguration) map[string]*DomainConfiguration {
	result := make(map[string]*DomainConfiguration)
	for name, value := range parsed {
		bc := blockConfigFromYaml(name, value.BlockConfiguration, BlockTypeDomain, types.CategoryDomain)
		result[name] = &DomainConfiguration{BlockConfiguration: *bc, DnsName: value.DnsName}
	}
	return result
}
