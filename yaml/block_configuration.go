package yaml

import "gopkg.in/nullstone-io/go-api-client.v0/types"

const RedactedValue = "••••••••••"

type BlockConfiguration struct {
	ModuleSource        string            `yaml:"module" json:"module"`
	ModuleSourceVersion *string           `yaml:"module_version,omitempty" json:"moduleVersion"`
	Variables           map[string]any    `yaml:"vars" json:"vars"`
	Connections         ConnectionTargets `yaml:"connections" json:"connections"`
	IsShared            bool              `yaml:"is_shared" json:"isShared"`
}

func BlockConfigurationFromWorkspaceConfig(stackId, envId int64, config types.WorkspaceConfig) BlockConfiguration {
	return BlockConfiguration{
		ModuleSource:        config.Source,
		ModuleSourceVersion: &config.SourceVersion,
		Variables:           VariablesFromWorkspaceConfig(config.Variables),
		Connections:         ConnectionsFromWorkspaceConfig(stackId, envId, config.Connections),
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

func ConnectionsFromWorkspaceConfig(stackId, envId int64, connections types.Connections) ConnectionTargets {
	result := ConnectionTargets{}
	for name, conn := range connections {
		if conn.Reference == nil {
			continue
		}
		target := ConnectionTarget{}
		// we only need to populate the stack name if it is different then the root workspace
		if conn.Reference.StackId != 0 && conn.Reference.StackId != stackId {
			target.StackName = conn.Reference.StackName
		}
		target.BlockName = conn.Reference.BlockName
		// we only need to populate the env name if it is different then the root workspace
		if conn.Reference.EnvId != nil && *conn.Reference.EnvId != envId {
			target.EnvName = conn.Reference.EnvName
		}
		result[name] = target
	}
	return result
}
