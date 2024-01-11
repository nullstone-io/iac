package config

import (
	"fmt"
	"github.com/BSick7/go-api/errors"
	"github.com/nullstone-io/iac/core"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
)

type IngressConfiguration struct {
	Name                string                 `yaml:"-" json:"name"`
	ModuleSource        string                 `yaml:"module" json:"module"`
	ModuleSourceVersion *string                `yaml:"module_version,omitempty" json:"moduleVersion"`
	Variables           map[string]any         `yaml:"vars" json:"vars"`
	Connections         core.ConnectionTargets `yaml:"connections" json:"connections"`
}

func (i IngressConfiguration) Validate(resolver *find.ResourceResolver, configBlocks []core.BlockConfiguration) (errors.ValidationErrors, error) {
	yamlPath := fmt.Sprintf("ingresses.%s", i.Name)
	return ValidateBlock(resolver, configBlocks, yamlPath, "ingress/*/*", i.ModuleSource, *i.ModuleSourceVersion, i.Variables, i.Connections, nil)
}

func (i *IngressConfiguration) Normalize(resolver *find.ResourceResolver) error {
	return core.NormalizeConnectionTargets(i.Connections, resolver)
}
