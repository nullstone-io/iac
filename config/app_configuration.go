package config

import (
	"fmt"
	"github.com/nullstone-io/iac/core"
	"github.com/nullstone-io/iac/yaml"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
)

type AppConfiguration struct {
	BlockConfiguration

	EnvVariables map[string]string
	Capabilities CapabilityConfigurations
}

func convertCapabilities(parsed yaml.CapabilityConfigurations) []CapabilityConfiguration {
	result := make([]CapabilityConfiguration, len(parsed))
	for i, capValue := range parsed {
		moduleVersion := "latest"
		if capValue.ModuleSourceVersion != nil {
			moduleVersion = *capValue.ModuleSourceVersion
		}
		result[i] = CapabilityConfiguration{
			ModuleSource:        capValue.ModuleSource,
			ModuleSourceVersion: moduleVersion,
			Variables:           capValue.Variables,
			Connections:         convertConnections(capValue.Connections),
			Namespace:           capValue.Namespace,
		}
	}
	return result
}

func convertAppConfigurations(parsed map[string]yaml.AppConfiguration) map[string]AppConfiguration {
	apps := make(map[string]AppConfiguration)
	for appName, appValue := range parsed {
		// set a default module version if not provided
		moduleVersion := "latest"
		if appValue.ModuleSourceVersion != nil {
			moduleVersion = *appValue.ModuleSourceVersion
		}
		app := AppConfiguration{
			BlockConfiguration: BlockConfiguration{
				Type:                BlockTypeApplication,
				Name:                appName,
				ModuleSource:        appValue.ModuleSource,
				ModuleSourceVersion: moduleVersion,
				Variables:           appValue.Variables,
				Connections:         convertConnections(appValue.Connections),
			},
			EnvVariables: appValue.EnvVariables,
			Capabilities: convertCapabilities(appValue.Capabilities),
		}
		apps[appName] = app
	}
	return apps
}

func (a AppConfiguration) Validate(resolver *find.ResourceResolver, configBlocks []BlockConfiguration) error {
	yamlPath := fmt.Sprintf("apps.%s", a.Name)
	contract := fmt.Sprintf("app/*/*")
	return ValidateBlock(resolver, configBlocks, yamlPath, contract, a.ModuleSource, a.ModuleSourceVersion, a.Variables, a.Connections, a.Capabilities)
}

func (a *AppConfiguration) Normalize(resolver *find.ResourceResolver) error {
	err := core.NormalizeConnectionTargets(a.Connections, resolver)
	if err != nil {
		return err
	}
	return a.Capabilities.Normalize(resolver)
}
