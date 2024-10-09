package config

import (
	"context"
	"fmt"
	"github.com/BSick7/go-api/errors"
	"github.com/nullstone-io/iac/core"
	"github.com/nullstone-io/iac/yaml"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
)

type DatastoreConfiguration struct {
	BlockConfiguration
}

func convertDatastoreConfigurations(parsed map[string]yaml.DatastoreConfiguration) map[string]DatastoreConfiguration {
	result := make(map[string]DatastoreConfiguration)
	for datastoreName, datastoreValue := range parsed {
		ds := DatastoreConfiguration{
			BlockConfiguration: blockConfigFromYaml(datastoreName, datastoreValue.BlockConfiguration, BlockTypeDatastore),
		}
		result[datastoreName] = ds
	}
	return result
}

func (d *DatastoreConfiguration) Validate(ctx context.Context, resolver *find.ResourceResolver, ic core.IacContext) errors.ValidationErrors {
	pc := core.NewYamlPathContext("datastores", d.Name)
	contract := fmt.Sprintf("datastore/*/*")
	return ValidateBlock(ctx, resolver, ic, pc, contract, d.ModuleSource, d.ModuleSourceVersion, d.Variables, d.Connections, nil, nil)
}
