package overrides

import "github.com/nullstone-io/iac/yaml"

type NetworkOverrides struct {
	BlockOverrides
}

func convertNetworkOverrides(parsed map[string]yaml.NetworkOverrides) map[string]NetworkOverrides {
	result := make(map[string]NetworkOverrides)
	for networkName, networkValue := range parsed {
		n := NetworkOverrides{
			BlockOverrides: BlockOverrides{
				Name:      networkName,
				Variables: networkValue.Variables,
			},
		}
		result[networkName] = n
	}
	return result
}
