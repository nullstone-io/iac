package overrides

import "github.com/nullstone-io/iac/yaml"

type ClusterNamespaceOverrides struct {
	BlockOverrides
}

func convertClusterNamespaceOverrides(parsed map[string]yaml.ClusterNamespaceOverrides) map[string]ClusterNamespaceOverrides {
	result := make(map[string]ClusterNamespaceOverrides)
	for clusterNamespaceName, clusterNamespaceValue := range parsed {
		cn := ClusterNamespaceOverrides{
			BlockOverrides: BlockOverrides{
				Name:      clusterNamespaceName,
				Variables: clusterNamespaceValue.Variables,
			},
		}
		result[clusterNamespaceName] = cn
	}
	return result
}
