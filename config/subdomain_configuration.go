package config

import (
	"fmt"
	"github.com/nullstone-io/iac/yaml"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
)

type SubdomainConfiguration struct {
	BlockConfiguration

	DnsName string
}

func convertSubdomainConfigurations(parsed map[string]yaml.SubdomainConfiguration) map[string]SubdomainConfiguration {
	result := make(map[string]SubdomainConfiguration)
	for subName, subValue := range parsed {
		// set a default module version if not provided
		moduleVersion := "latest"
		if subValue.ModuleSourceVersion != nil {
			moduleVersion = *subValue.ModuleSourceVersion
		}
		sub := SubdomainConfiguration{
			BlockConfiguration: BlockConfiguration{
				Type:                BlockTypeSubdomain,
				Name:                subName,
				ModuleSource:        subValue.ModuleSource,
				ModuleSourceVersion: moduleVersion,
				Variables:           subValue.Variables,
				Connections:         convertConnections(subValue.Connections),
			},
			DnsName: subValue.DnsName,
		}
		result[subName] = sub
	}
	return result
}

func (s SubdomainConfiguration) Validate(resolver *find.ResourceResolver) error {
	yamlPath := fmt.Sprintf("subdomains.%s", s.Name)
	contract := fmt.Sprintf("subdomain/*/*")
	return ValidateBlock(resolver, yamlPath, contract, s.ModuleSource, s.ModuleSourceVersion, s.Variables, s.Connections, nil, nil)
}
