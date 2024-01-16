package overrides

import (
	"github.com/nullstone-io/iac/yaml"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

type AppOverrides struct {
	BlockOverrides

	EnvVariables map[string]string
	Capabilities CapabilityOverrides
}

func convertConnections(parsed map[string]yaml.ConnectionTarget) map[string]types.ConnectionTarget {
	result := make(map[string]types.ConnectionTarget)
	for key, conn := range parsed {
		result[key] = types.ConnectionTarget{
			StackId:   conn.StackId,
			StackName: conn.StackName,
			BlockId:   conn.BlockId,
			BlockName: conn.BlockName,
			EnvId:     conn.EnvId,
			EnvName:   conn.EnvName,
		}
	}
	return result
}

func convertCapabilities(parsed yaml.CapabilityOverrides) CapabilityOverrides {
	result := make(CapabilityOverrides, len(parsed))
	for i, capValue := range parsed {
		moduleVersion := "latest"
		if capValue.ModuleSourceVersion != nil {
			moduleVersion = *capValue.ModuleSourceVersion
		}
		result[i] = CapabilityOverride{
			ModuleSource:        capValue.ModuleSource,
			ModuleSourceVersion: moduleVersion,
			Variables:           capValue.Variables,
			Connections:         convertConnections(capValue.Connections),
			Namespace:           capValue.Namespace,
		}
	}
	return result
}

func convertAppOverrides(parsed map[string]yaml.AppOverrides) map[string]AppOverrides {
	apps := make(map[string]AppOverrides)
	for appName, appValue := range parsed {
		app := AppOverrides{
			BlockOverrides: BlockOverrides{
				Name:      appName,
				Variables: appValue.Variables,
			},
			EnvVariables: appValue.EnvVariables,
			Capabilities: convertCapabilities(appValue.Capabilities),
		}
		apps[appName] = app
	}
	return apps
}

func (a *AppOverrides) Normalize(resolver *find.ResourceResolver) error {
	return a.Capabilities.Normalize(resolver)
}
