package config

import (
	"context"
	"fmt"
	"github.com/nullstone-io/iac/yaml"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
)

type DomainConfiguration struct {
	BlockConfiguration

	DnsName string `json:"dnsName"`
}

func convertDomainConfigurations(parsed map[string]yaml.DomainConfiguration) map[string]DomainConfiguration {
	result := make(map[string]DomainConfiguration)
	for domainName, domainValue := range parsed {
		domain := DomainConfiguration{
			BlockConfiguration: blockConfigFromYaml(domainName, domainValue.BlockConfiguration, BlockTypeDomain),
			DnsName:            domainValue.DnsName,
		}
		result[domainName] = domain
	}
	return result
}

func (d DomainConfiguration) Validate(ctx context.Context, resolver *find.ResourceResolver, repoName, filename string) error {
	yamlPath := fmt.Sprintf("domains.%s", d.Name)
	contract := fmt.Sprintf("domain/*/*")
	return ValidateBlock(ctx, resolver, repoName, filename, yamlPath, contract, d.ModuleSource, d.ModuleSourceVersion, d.Variables, d.Connections, nil, nil)
}
