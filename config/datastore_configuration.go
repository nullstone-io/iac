package config

import (
	"github.com/nullstone-io/iac/yaml"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

type DatastoreConfiguration struct {
	BlockConfiguration
}

func convertDatastoreConfigurations(parsed map[string]yaml.DatastoreConfiguration) map[string]*DatastoreConfiguration {
	result := make(map[string]*DatastoreConfiguration)
	for name, value := range parsed {
		bc := blockConfigFromYaml(name, value.BlockConfiguration, BlockTypeDatastore, types.CategoryDatastore)
		result[name] = &DatastoreConfiguration{BlockConfiguration: *bc}
	}
	return result
}
