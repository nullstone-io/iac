package config

import (
	"github.com/nullstone-io/iac/yaml"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

type ClusterConfiguration struct {
	BlockConfiguration
}

func convertClusterConfigurations(parsed map[string]yaml.ClusterConfiguration) map[string]*ClusterConfiguration {
	result := make(map[string]*ClusterConfiguration)
	for name, value := range parsed {
		bc := blockConfigFromYaml(name, value.BlockConfiguration, BlockTypeCluster, types.CategoryCluster)
		result[name] = &ClusterConfiguration{BlockConfiguration: *bc}
	}
	return result
}
