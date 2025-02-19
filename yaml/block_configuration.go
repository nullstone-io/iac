package yaml

import "gopkg.in/nullstone-io/go-api-client.v0/types"

const RedactedValue = "••••••••••"

type BlockConfiguration struct {
	ModuleSource     string                `yaml:"module" json:"module"`
	ModuleConstraint *string               `yaml:"module_version,omitempty" json:"moduleVersion"`
	Variables        map[string]any        `yaml:"vars,omitempty" json:"vars"`
	Connections      ConnectionConstraints `yaml:"connections,omitempty" json:"connections"`
	IsShared         bool                  `yaml:"is_shared,omitempty" json:"isShared"`
}

func BlockConfigurationFromWorkspaceConfig(stackId, envId int64, config types.WorkspaceConfig) BlockConfiguration {
	return BlockConfiguration{
		ModuleSource:     config.Source,
		ModuleConstraint: &config.SourceConstraint,
		Variables:        VariablesFromWorkspaceConfig(config.Variables),
		Connections:      ConnectionsFromWorkspaceConfig(stackId, envId, config.Connections),
	}
}

func VariablesFromWorkspaceConfig(variables types.Variables) map[string]any {
	result := map[string]any{}
	for name, value := range variables {
		if value.HasValue() {
			if value.Sensitive {
				result[name] = RedactedValue
			} else {
				result[name] = value.Value
			}
		}
	}
	return result
}

func ConnectionsFromWorkspaceConfig(stackId, envId int64, connections types.Connections) ConnectionConstraints {
	result := ConnectionConstraints{}
	for name, conn := range connections {
		if conn.DesiredTarget == nil {
			continue
		}
		result[name] = ConnectionConstraint{
			StackName: conn.DesiredTarget.StackName,
			EnvName:   conn.DesiredTarget.EnvName,
			BlockName: conn.DesiredTarget.BlockName,
		}
	}
	return result
}
