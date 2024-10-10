package config

import (
	"github.com/nullstone-io/iac/yaml"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

type DomainConfiguration struct {
	BlockConfiguration

	DnsName string `json:"dnsName"`
}

func convertDomainConfigurations(parsed map[string]yaml.DomainConfiguration) map[string]DomainConfiguration {
	result := make(map[string]DomainConfiguration)
	for domainName, domainValue := range parsed {
		domain := DomainConfiguration{
			BlockConfiguration: blockConfigFromYaml(domainName, domainValue.BlockConfiguration, BlockTypeDomain, types.CategoryDomain),
			DnsName:            domainValue.DnsName,
		}
		result[domainName] = domain
	}
	return result
}
