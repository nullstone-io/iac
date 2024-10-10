package config

import (
	"github.com/nullstone-io/iac/yaml"
)

type IngressConfiguration struct {
	BlockConfiguration
}

func convertIngressConfigurations(parsed map[string]yaml.IngressConfiguration) map[string]IngressConfiguration {
	result := make(map[string]IngressConfiguration)
	for ingressName, ingressValue := range parsed {
		ingress := IngressConfiguration{
			BlockConfiguration: blockConfigFromYaml(ingressName, ingressValue.BlockConfiguration, BlockTypeIngress, "ingress"),
		}
		result[ingressName] = ingress
	}
	return result
}
