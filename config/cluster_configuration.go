package config

import (
	"fmt"
	"github.com/BSick7/go-api/errors"
	"github.com/nullstone-io/iac/core"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
)

type ClusterConfiguration struct {
	Name                string                 `yaml:"-" json:"name"`
	ModuleSource        string                 `yaml:"module" json:"module"`
	ModuleSourceVersion *string                `yaml:"module_version,omitempty" json:"moduleVersion"`
	Variables           map[string]any         `yaml:"vars" json:"vars"`
	Connections         core.ConnectionTargets `yaml:"connections" json:"connections"`
}

func (c ClusterConfiguration) Validate(resolver *find.ResourceResolver, configBlocks []core.BlockConfiguration) (errors.ValidationErrors, error) {
	yamlPath := fmt.Sprintf("clusters.%s", c.Name)
	return ValidateBlock(resolver, configBlocks, yamlPath, "cluster/*/*", c.ModuleSource, *c.ModuleSourceVersion, c.Variables, c.Connections, nil)
}

func (c *ClusterConfiguration) Normalize(resolver *find.ResourceResolver) error {
	return core.NormalizeConnectionTargets(c.Connections, resolver)
}
