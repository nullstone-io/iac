package config

import (
	"context"
	"fmt"
	"github.com/nullstone-io/iac/yaml"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
)

type DatastoreConfiguration struct {
	BlockConfiguration
}

func convertDatastoreConfigurations(parsed map[string]yaml.DatastoreConfiguration) map[string]DatastoreConfiguration {
	result := make(map[string]DatastoreConfiguration)
	for datastoreName, datastoreValue := range parsed {
		// set a default module version if not provided
		moduleVersion := "latest"
		if datastoreValue.ModuleSourceVersion != nil {
			moduleVersion = *datastoreValue.ModuleSourceVersion
		}
		ds := DatastoreConfiguration{
			BlockConfiguration: BlockConfiguration{
				Type:                BlockTypeDatastore,
				Name:                datastoreName,
				ModuleSource:        datastoreValue.ModuleSource,
				ModuleSourceVersion: moduleVersion,
				Variables:           datastoreValue.Variables,
				Connections:         convertConnections(datastoreValue.Connections),
				IsShared:            datastoreValue.IsShared,
			},
		}
		result[datastoreName] = ds
	}
	return result
}

func (d DatastoreConfiguration) Validate(ctx context.Context, resolver *find.ResourceResolver, repoName, filename string) error {
	yamlPath := fmt.Sprintf("datastores.%s", d.Name)
	contract := fmt.Sprintf("datastore/*/*")
	return ValidateBlock(ctx, resolver, repoName, filename, yamlPath, contract, d.ModuleSource, d.ModuleSourceVersion, d.Variables, d.Connections, nil, nil)
}
