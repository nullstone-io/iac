package overrides

import "github.com/nullstone-io/iac/yaml"

type ClusterOverrides struct {
	BlockOverrides
}

func convertClusterOverrides(parsed map[string]yaml.ClusterOverrides) map[string]ClusterOverrides {
	result := make(map[string]ClusterOverrides)
	for clusterName, clusterValue := range parsed {
		c := ClusterOverrides{
			BlockOverrides: BlockOverrides{
				Name:      clusterName,
				Variables: clusterValue.Variables,
			},
		}
		result[clusterName] = c
	}
	return result
}
