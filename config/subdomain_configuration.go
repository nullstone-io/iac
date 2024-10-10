package config

import (
	"github.com/nullstone-io/iac/yaml"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

type SubdomainConfiguration struct {
	BlockConfiguration

	DnsName string `json:"dnsName"`
}

func convertSubdomainConfigurations(parsed map[string]yaml.SubdomainConfiguration) map[string]SubdomainConfiguration {
	result := make(map[string]SubdomainConfiguration)
	for subName, subValue := range parsed {
		sub := SubdomainConfiguration{
			BlockConfiguration: blockConfigFromYaml(subName, subValue.BlockConfiguration, BlockTypeSubdomain, types.CategorySubdomain),
			DnsName:            subValue.DnsName,
		}
		result[subName] = sub
	}
	return result
}
