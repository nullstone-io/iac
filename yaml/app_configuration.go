package yaml

import "gopkg.in/nullstone-io/go-api-client.v0/types"

type CapabilityConfigurations []CapabilityConfiguration

type CapabilityConfiguration struct {
	ModuleSource        string            `yaml:"module,omitempty" json:"module"`
	ModuleSourceVersion *string           `yaml:"module_version,omitempty" json:"moduleVersion"`
	Variables           map[string]any    `yaml:"vars" json:"vars"`
	Connections         ConnectionTargets `yaml:"connections,omitempty" json:"connections"`
	Namespace           *string           `yaml:"namespace,omitempty" json:"namespace"`
}

type AppConfiguration struct {
	BlockConfiguration `yaml:",inline" json:",inline"`

	EnvVariables map[string]string        `yaml:"environment" json:"envVars"`
	Capabilities CapabilityConfigurations `yaml:"capabilities" json:"capabilities"`
}

func AppConfigurationFromWorkspaceConfig(stackId, envId int64, config types.WorkspaceConfig) AppConfiguration {
	envVars := map[string]string{}
	for name, v := range config.EnvVariables {
		if v.Sensitive {
			envVars[name] = RedactedValue
		} else {
			envVars[name] = v.Value
		}
	}
	caps := CapabilityConfigurations{}
	for _, v := range config.Capabilities {
		var namespace *string
		if v.Namespace != "" {
			namespace = &v.Namespace
		}
		caps = append(caps, CapabilityConfiguration{
			ModuleSource:        v.Source,
			ModuleSourceVersion: &v.SourceVersion,
			Variables:           VariablesFromWorkspaceConfig(v.Variables),
			Connections:         ConnectionsFromWorkspaceConfig(stackId, envId, v.Connections),
			Namespace:           namespace,
		})
	}
	return AppConfiguration{
		BlockConfiguration: BlockConfigurationFromWorkspaceConfig(stackId, envId, config),
		EnvVariables:       envVars,
		Capabilities:       caps,
	}
}
