package config

import (
	"github.com/nullstone-io/iac/yaml"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

type ClusterNamespaceConfiguration struct {
	BlockConfiguration
}

func convertClusterNamespaceConfigurations(parsed map[string]yaml.ClusterNamespaceConfiguration) map[string]ClusterNamespaceConfiguration {
	result := make(map[string]ClusterNamespaceConfiguration)
	for clusterNamespaceName, clusterNamespaceValue := range parsed {
		cn := ClusterNamespaceConfiguration{
			BlockConfiguration: blockConfigFromYaml(clusterNamespaceName, clusterNamespaceValue.BlockConfiguration, BlockTypeClusterNamespace, types.CategoryClusterNamespace),
		}
		result[clusterNamespaceName] = cn
	}
	return result
}
