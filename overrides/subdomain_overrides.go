package overrides

import "github.com/nullstone-io/iac/yaml"

type SubdomainOverrides struct {
	BlockOverrides
}

func convertSubdomainOverrides(parsed map[string]yaml.SubdomainOverrides) map[string]SubdomainOverrides {
	result := make(map[string]SubdomainOverrides)
	for subdomainName, subdomainValue := range parsed {
		s := SubdomainOverrides{
			BlockOverrides: BlockOverrides{
				Name:      subdomainName,
				Variables: subdomainValue.Variables,
			},
		}
		result[subdomainName] = s
	}
	return result
}
