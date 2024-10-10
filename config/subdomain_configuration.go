package config

import (
	"github.com/nullstone-io/iac/yaml"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

type SubdomainConfiguration struct {
	BlockConfiguration

	DnsName string `json:"dnsName"`
}

func convertSubdomainConfigurations(parsed map[string]yaml.SubdomainConfiguration) map[string]*SubdomainConfiguration {
	result := make(map[string]*SubdomainConfiguration)
	for name, value := range parsed {
		bc := blockConfigFromYaml(name, value.BlockConfiguration, BlockTypeSubdomain, types.CategorySubdomain)
		result[name] = &SubdomainConfiguration{BlockConfiguration: *bc, DnsName: value.DnsName}
	}
	return result
}
