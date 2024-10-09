package config

import (
	"context"
	"fmt"
	"github.com/BSick7/go-api/errors"
	"github.com/nullstone-io/iac/core"
	"github.com/nullstone-io/iac/yaml"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
)

type ClusterNamespaceConfiguration struct {
	BlockConfiguration
}

func convertClusterNamespaceConfigurations(parsed map[string]yaml.ClusterNamespaceConfiguration) map[string]ClusterNamespaceConfiguration {
	result := make(map[string]ClusterNamespaceConfiguration)
	for clusterNamespaceName, clusterNamespaceValue := range parsed {
		cn := ClusterNamespaceConfiguration{
			BlockConfiguration: blockConfigFromYaml(clusterNamespaceName, clusterNamespaceValue.BlockConfiguration, BlockTypeClusterNamespace),
		}
		result[clusterNamespaceName] = cn
	}
	return result
}

func (cn ClusterNamespaceConfiguration) Validate(ctx context.Context, resolver *find.ResourceResolver, ic core.IacContext) errors.ValidationErrors {
	pc := core.NewYamlPathContext("cluster_namespaces", cn.Name)
	contract := fmt.Sprintf("cluster-namespace/*/*")
	return ValidateBlock(ctx, resolver, ic, pc, contract, cn.ModuleSource, cn.ModuleSourceVersion, cn.Variables, cn.Connections, nil, nil)
}
