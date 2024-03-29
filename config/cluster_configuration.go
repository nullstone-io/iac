package config

import (
	"context"
	"fmt"
	"github.com/nullstone-io/iac/yaml"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
)

type ClusterConfiguration struct {
	BlockConfiguration
}

func convertClusterConfigurations(parsed map[string]yaml.ClusterConfiguration) map[string]ClusterConfiguration {
	result := make(map[string]ClusterConfiguration)
	for clusterName, clusterValue := range parsed {
		// set a default module version if not provided
		moduleVersion := "latest"
		if clusterValue.ModuleSourceVersion != nil {
			moduleVersion = *clusterValue.ModuleSourceVersion
		}
		cluster := ClusterConfiguration{
			BlockConfiguration: BlockConfiguration{
				Type:                BlockTypeCluster,
				Name:                clusterName,
				ModuleSource:        clusterValue.ModuleSource,
				ModuleSourceVersion: moduleVersion,
				Variables:           clusterValue.Variables,
				Connections:         convertConnections(clusterValue.Connections),
				IsShared:            clusterValue.IsShared,
			},
		}
		result[clusterName] = cluster
	}
	return result
}

func (c ClusterConfiguration) Validate(ctx context.Context, resolver *find.ResourceResolver, repoName, filename string) error {
	yamlPath := fmt.Sprintf("clusters.%s", c.Name)
	contract := fmt.Sprintf("cluster/*/*")
	return ValidateBlock(ctx, resolver, repoName, filename, yamlPath, contract, c.ModuleSource, c.ModuleSourceVersion, c.Variables, c.Connections, nil, nil)
}
