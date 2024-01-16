package overrides

import "github.com/nullstone-io/iac/yaml"

type DatastoreOverrides struct {
	BlockOverrides
}

func convertDatastoreOverrides(parsed map[string]yaml.DatastoreOverrides) map[string]DatastoreOverrides {
	result := make(map[string]DatastoreOverrides)
	for datastoreName, datastoreValue := range parsed {
		d := DatastoreOverrides{
			BlockOverrides: BlockOverrides{
				Name:      datastoreName,
				Variables: datastoreValue.Variables,
			},
		}
		result[datastoreName] = d
	}
	return result
}
