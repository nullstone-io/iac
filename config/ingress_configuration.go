package config

import (
	"fmt"
	"github.com/nullstone-io/iac/yaml"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
)

type IngressConfiguration struct {
	BlockConfiguration
}

func convertIngressConfigurations(parsed map[string]yaml.IngressConfiguration) map[string]IngressConfiguration {
	result := make(map[string]IngressConfiguration)
	for ingressName, ingressValue := range parsed {
		// set a default module version if not provided
		moduleVersion := "latest"
		if ingressValue.ModuleSourceVersion != nil {
			moduleVersion = *ingressValue.ModuleSourceVersion
		}
		ingress := IngressConfiguration{
			BlockConfiguration: BlockConfiguration{
				Type:                BlockTypeIngress,
				Name:                ingressName,
				ModuleSource:        ingressValue.ModuleSource,
				ModuleSourceVersion: moduleVersion,
				Variables:           ingressValue.Variables,
				Connections:         convertConnections(ingressValue.Connections),
			},
		}
		result[ingressName] = ingress
	}
	return result
}

func (i IngressConfiguration) Validate(resolver *find.ResourceResolver) error {
	yamlPath := fmt.Sprintf("ingresses.%s", i.Name)
	contract := fmt.Sprintf("ingress/*/*")
	return ValidateBlock(resolver, yamlPath, contract, i.ModuleSource, i.ModuleSourceVersion, i.Variables, i.Connections, nil, nil)
}
