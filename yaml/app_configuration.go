package yaml

import "gopkg.in/nullstone-io/go-api-client.v0/types"

type AppConfiguration struct {
	BlockConfiguration `yaml:",inline" json:",inline"`

	Framework    string                   `yaml:"framework,omitempty" json:"framework"`
	EnvVariables map[string]string        `yaml:"environment,omitempty" json:"envVars"`
	Capabilities CapabilityConfigurations `yaml:"capabilities,omitempty" json:"capabilities"`
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
			Name:                v.Name,
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
