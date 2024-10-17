package config

import (
	"github.com/nullstone-io/iac/yaml"
)

type IngressConfiguration struct {
	BlockConfiguration
}

func convertIngressConfigurations(parsed map[string]yaml.IngressConfiguration) map[string]*IngressConfiguration {
	result := make(map[string]*IngressConfiguration)
	for name, value := range parsed {
		bc := blockConfigFromYaml(name, value.BlockConfiguration, BlockTypeIngress, "ingress")
		result[name] = &IngressConfiguration{BlockConfiguration: *bc}
	}
	return result
}
