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

func AppConfigurationFromWorkspaceConfig(config types.WorkspaceConfig) AppConfiguration {
	envVars := map[string]string{}
	for k, v := range config.EnvVariables {
		envVars[k] = v.Value
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
			Connections:         ConnectionsFromWorkspaceConfig(v.Connections),
			Namespace:           namespace,
		})
	}
	return AppConfiguration{
		BlockConfiguration: BlockConfigurationFromWorkspaceConfig(config),
		EnvVariables:       envVars,
		Capabilities:       caps,
	}
}
