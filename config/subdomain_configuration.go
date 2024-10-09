package config

import (
	"context"
	"fmt"
	"github.com/BSick7/go-api/errors"
	"github.com/nullstone-io/iac/core"
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

func (s SubdomainConfiguration) Validate(ctx context.Context, resolver *find.ResourceResolver, ic core.IacContext) errors.ValidationErrors {
	pc := core.NewYamlPathContext("subdomains", s.Name)
	contract := fmt.Sprintf("subdomain/*/*")
	return ValidateBlock(ctx, resolver, ic, pc, contract, s.ModuleSource, s.ModuleSourceVersion, s.Variables, s.Connections, nil, nil)
}
