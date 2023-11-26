package config

import (
	"fmt"
	"github.com/BSick7/go-api/errors"
	"github.com/nullstone-io/iac/core"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
)

type DatastoreConfiguration struct {
	Name                string                 `yaml:"-" json:"name"`
	ModuleSource        string                 `yaml:"module" json:"module"`
	ModuleSourceVersion *string                `yaml:"module_version,omitempty" json:"moduleVersion"`
	Variables           map[string]any         `yaml:"vars,omitempty" json:"vars"`
	Connections         core.ConnectionTargets `yaml:"connections,omitempty" json:"connections"`
}

func (d DatastoreConfiguration) Validate(resolver *find.ResourceResolver, configBlocks []core.BlockConfiguration) (errors.ValidationErrors, error) {
	yamlPath := fmt.Sprintf("datastores.%s", d.Name)
	return ValidateBlock(resolver, configBlocks, yamlPath, "datastore/*/*", d.ModuleSource, *d.ModuleSourceVersion, d.Variables, d.Connections, nil)
}

func (d *DatastoreConfiguration) Normalize(resolver *find.ResourceResolver) error {
	return core.NormalizeConnectionTargets(d.Connections, resolver)
}
