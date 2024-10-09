package config

import (
	"context"
	"fmt"
	"github.com/BSick7/go-api/errors"
	"github.com/nullstone-io/iac/core"
	"github.com/nullstone-io/iac/yaml"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
)

type AppConfiguration struct {
	BlockConfiguration

	EnvVariables map[string]string        `json:"envVars"`
	Capabilities CapabilityConfigurations `json:"capabilities"`
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
		app := AppConfiguration{
			BlockConfiguration: blockConfigFromYaml(appName, appValue.BlockConfiguration, BlockTypeApplication),
			EnvVariables:       appValue.EnvVariables,
			Capabilities:       convertCapabilities(appValue.Capabilities),
		}
		apps[appName] = app
	}
	return apps
}

func (a *AppConfiguration) Validate(ctx context.Context, resolver *find.ResourceResolver, ic core.IacContext) errors.ValidationErrors {
	pc := core.NewYamlPathContext("apps", a.Name)
	contract := fmt.Sprintf("app/*/*")
	return ValidateBlock(ctx, resolver, ic, pc, contract, a.ModuleSource, a.ModuleSourceVersion, a.Variables, a.Connections, a.EnvVariables, a.Capabilities)
}

func (a *AppConfiguration) Normalize(ctx context.Context, resolver *find.ResourceResolver) error {
	err := NormalizeConnectionTargets(ctx, a.Connections, resolver)
	if err != nil {
		return err
	}
	return a.Capabilities.Normalize(ctx, resolver)
}
