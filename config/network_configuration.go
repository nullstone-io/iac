package config

import (
	"github.com/nullstone-io/iac/yaml"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

type NetworkConfiguration struct {
	BlockConfiguration
}

func convertNetworkConfigurations(parsed map[string]yaml.NetworkConfiguration) map[string]*NetworkConfiguration {
	result := make(map[string]*NetworkConfiguration)
	for networkName, networkValue := range parsed {
		bc := blockConfigFromYaml(networkName, networkValue.BlockConfiguration, BlockTypeNetwork, types.CategoryNetwork)
		result[networkName] = &NetworkConfiguration{BlockConfiguration: *bc}
	}
	return result
}
