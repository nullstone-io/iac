package overrides

import "github.com/nullstone-io/iac/yaml"

type IngressOverrides struct {
	BlockOverrides
}

func convertIngressOverrides(parsed map[string]yaml.IngressOverrides) map[string]IngressOverrides {
	result := make(map[string]IngressOverrides)
	for ingressName, ingressValue := range parsed {
		i := IngressOverrides{
			BlockOverrides: BlockOverrides{
				Name:      ingressName,
				Variables: ingressValue.Variables,
			},
		}
		result[ingressName] = i
	}
	return result
}
