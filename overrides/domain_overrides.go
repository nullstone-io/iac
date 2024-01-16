package overrides

import "github.com/nullstone-io/iac/yaml"

type DomainOverrides struct {
	BlockOverrides
}

func convertDomainOverrides(parsed map[string]yaml.DomainOverrides) map[string]DomainOverrides {
	result := make(map[string]DomainOverrides)
	for domainName, domainValue := range parsed {
		d := DomainOverrides{
			BlockOverrides: BlockOverrides{
				Name:      domainName,
				Variables: domainValue.Variables,
			},
		}
		result[domainName] = d
	}
	return result
}
