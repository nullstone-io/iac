package config

import (
	"fmt"
	"github.com/BSick7/go-api/errors"
	"github.com/nullstone-io/iac/core"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
)

type ClusterNamespaceConfiguration struct {
	Name                string                 `yaml:"-" json:"name"`
	ModuleSource        string                 `yaml:"module" json:"module"`
	ModuleSourceVersion *string                `yaml:"module_version,omitempty" json:"moduleVersion"`
	Variables           map[string]any         `yaml:"vars" json:"vars"`
	Connections         core.ConnectionTargets `yaml:"connections" json:"connections"`
}

func (cn ClusterNamespaceConfiguration) Validate(resolver *find.ResourceResolver, configBlocks []core.BlockConfiguration) (errors.ValidationErrors, error) {
	yamlPath := fmt.Sprintf("cluster_namespaces.%s", cn.Name)
	return ValidateBlock(resolver, configBlocks, yamlPath, "*/*/*", cn.ModuleSource, *cn.ModuleSourceVersion, cn.Variables, cn.Connections, nil)
}

func (cn *ClusterNamespaceConfiguration) Normalize(resolver *find.ResourceResolver) error {
	return core.NormalizeConnectionTargets(cn.Connections, resolver)
}
