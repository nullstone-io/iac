package yaml

import "gopkg.in/nullstone-io/go-api-client.v0/types"

type BlockConfiguration struct {
	ModuleSource        string            `yaml:"module" json:"module"`
	ModuleSourceVersion *string           `yaml:"module_version,omitempty" json:"moduleVersion"`
	Variables           map[string]any    `yaml:"vars" json:"vars"`
	Connections         ConnectionTargets `yaml:"connections" json:"connections"`
	IsShared            bool              `yaml:"is_shared" json:"isShared"`
}

func BlockConfigurationFromWorkspaceConfig(config types.WorkspaceConfig) BlockConfiguration {
	return BlockConfiguration{
		ModuleSource:        config.Source,
		ModuleSourceVersion: &config.SourceVersion,
		Variables:           VariablesFromWorkspaceConfig(config.Variables),
		Connections:         ConnectionsFromWorkspaceConfig(config.Connections),
	}
}

func VariablesFromWorkspaceConfig(variables types.Variables) map[string]any {
	result := map[string]any{}
	for name, value := range variables {
		if value.HasValue() {
			result[name] = value.Value
		}
	}
	return result
}

func ConnectionsFromWorkspaceConfig(connections types.Connections) ConnectionTargets {
	result := ConnectionTargets{}
	for name, connection := range connections {
		if connection.Reference != nil {
			result[name] = ConnectionTarget(*connection.Reference)
		}
	}
	return result
}
