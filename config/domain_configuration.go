package config

import (
	"context"
	"fmt"
	"github.com/BSick7/go-api/errors"
	"github.com/nullstone-io/iac/core"
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

func (d *DomainConfiguration) Validate(ctx context.Context, resolver *find.ResourceResolver, ic core.IacContext) errors.ValidationErrors {
	pc := core.NewYamlPathContext("domains", d.Name)
	contract := fmt.Sprintf("domain/*/*")
	return ValidateBlock(ctx, resolver, ic, pc, contract, d.ModuleSource, d.ModuleSourceVersion, d.Variables, d.Connections, nil, nil)
}
