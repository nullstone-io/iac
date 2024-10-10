package config

import (
	"github.com/nullstone-io/iac/yaml"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

type ClusterConfiguration struct {
	BlockConfiguration
}

func convertClusterConfigurations(parsed map[string]yaml.ClusterConfiguration) map[string]ClusterConfiguration {
	result := make(map[string]ClusterConfiguration)
	for clusterName, clusterValue := range parsed {
		cluster := ClusterConfiguration{
			BlockConfiguration: blockConfigFromYaml(clusterName, clusterValue.BlockConfiguration, BlockTypeCluster, types.CategoryCluster),
		}
		result[clusterName] = cluster
	}
	return result
}
