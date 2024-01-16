package config

import (
	"fmt"
	"github.com/nullstone-io/iac/core"
	"github.com/nullstone-io/iac/yaml"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
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

type InvalidConfigurationError struct {
	Err error
}

func (e InvalidConfigurationError) Error() string {
	return fmt.Sprintf("invalid app configuration: %s", e.Err.Error())
}

func (a AppConfiguration) GetCapabilities(orgName string, stackId, blockId, envId int64) ([]types.Capability, error) {
	caps := make([]types.Capability, len(a.Capabilities))
	for i, cap := range a.Capabilities {
		updateCap := types.Capability{
			OrgName:             orgName,
			AppId:               blockId,
			EnvId:               envId,
			ModuleSource:        cap.ModuleSource,
			ModuleSourceVersion: cap.ModuleSourceVersion,
			Connections:         map[string]types.ConnectionTarget{},
		}
		if cap.Namespace != nil {
			updateCap.Namespace = *cap.Namespace
		}
		for key, conn := range cap.Connections {
			target := types.ConnectionTarget{}
			// each connection must have a block_name to identify which block it is connected to
			if conn.BlockName == "" {
				return nil, InvalidConfigurationError{fmt.Errorf("The connection (%s) must have a block_name to identify which block it is connected to.", key)}
			}
			target.BlockName = conn.BlockName
			// each connection must also have a stack_id
			if conn.StackId != 0 {
				target.StackId = conn.StackId
			} else {
				target.StackId = stackId
			}
			target.EnvId = conn.EnvId
			updateCap.Connections[key] = target
		}
		caps[i] = updateCap
	}
	return caps, nil
}
