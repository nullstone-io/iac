package config

import (
	"context"
	"fmt"
	"github.com/BSick7/go-api/errors"
	"github.com/nullstone-io/iac/core"
	"github.com/nullstone-io/iac/yaml"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
)

type IngressConfiguration struct {
	BlockConfiguration
}

func convertIngressConfigurations(parsed map[string]yaml.IngressConfiguration) map[string]IngressConfiguration {
	result := make(map[string]IngressConfiguration)
	for ingressName, ingressValue := range parsed {
		ingress := IngressConfiguration{
			BlockConfiguration: blockConfigFromYaml(ingressName, ingressValue.BlockConfiguration, BlockTypeIngress),
		}
		result[ingressName] = ingress
	}
	return result
}

func (i IngressConfiguration) Validate(ctx context.Context, resolver *find.ResourceResolver, ic core.IacContext) errors.ValidationErrors {
	pc := core.NewYamlPathContext("ingresses", i.Name)
	contract := fmt.Sprintf("ingress/*/*")
	return ValidateBlock(ctx, resolver, ic, pc, contract, i.ModuleSource, i.ModuleSourceVersion, i.Variables, i.Connections, nil, nil)
}
