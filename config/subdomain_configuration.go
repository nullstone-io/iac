package config

import (
	"context"
	"fmt"
	"github.com/nullstone-io/iac/yaml"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
)

type SubdomainConfiguration struct {
	BlockConfiguration

	DnsName string `json:"dnsName"`
}

func convertSubdomainConfigurations(parsed map[string]yaml.SubdomainConfiguration) map[string]SubdomainConfiguration {
	result := make(map[string]SubdomainConfiguration)
	for subName, subValue := range parsed {
		sub := SubdomainConfiguration{
			BlockConfiguration: blockConfigFromYaml(subName, subValue.BlockConfiguration, BlockTypeSubdomain),
			DnsName:            subValue.DnsName,
		}
		result[subName] = sub
	}
	return result
}

func (s SubdomainConfiguration) Validate(ctx context.Context, resolver *find.ResourceResolver, repoName, filename string) error {
	yamlPath := fmt.Sprintf("subdomains.%s", s.Name)
	contract := fmt.Sprintf("subdomain/*/*")
	return ValidateBlock(ctx, resolver, repoName, filename, yamlPath, contract, s.ModuleSource, s.ModuleSourceVersion, s.Variables, s.Connections, nil, nil)
}
