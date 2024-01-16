package config

import (
	"fmt"
	"github.com/nullstone-io/iac/yaml"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
)

type ClusterNamespaceConfiguration struct {
	BlockConfiguration
}

func convertClusterNamespaceConfigurations(parsed map[string]yaml.ClusterNamespaceConfiguration) map[string]ClusterNamespaceConfiguration {
	result := make(map[string]ClusterNamespaceConfiguration)
	for clusterNamespaceName, clusterNamespaceValue := range parsed {
		// set a default module version if not provided
		moduleVersion := "latest"
		if clusterNamespaceValue.ModuleSourceVersion != nil {
			moduleVersion = *clusterNamespaceValue.ModuleSourceVersion
		}
		cn := ClusterNamespaceConfiguration{
			BlockConfiguration: BlockConfiguration{
				Type:                BlockTypeClusterNamespace,
				Name:                clusterNamespaceName,
				ModuleSource:        clusterNamespaceValue.ModuleSource,
				ModuleSourceVersion: moduleVersion,
				Variables:           clusterNamespaceValue.Variables,
				Connections:         convertConnections(clusterNamespaceValue.Connections),
			},
		}
		result[clusterNamespaceName] = cn
	}
	return result
}

func (cn ClusterNamespaceConfiguration) Validate(resolver *find.ResourceResolver, configBlocks []BlockConfiguration) error {
	yamlPath := fmt.Sprintf("cluster_namespaces.%s", cn.Name)
	contract := fmt.Sprintf("cluster-namespace/*/*")
	return ValidateBlock(resolver, configBlocks, yamlPath, contract, cn.ModuleSource, cn.ModuleSourceVersion, cn.Variables, cn.Connections, nil)
}
