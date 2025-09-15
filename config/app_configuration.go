package config

import (
	"context"
	"strings"

	"github.com/nullstone-io/iac/core"
	"github.com/nullstone-io/iac/yaml"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

var (
	_ core.ChangeApplier = &AppConfiguration{}
)

type AppConfiguration struct {
	BlockConfiguration

	Framework    string                   `json:"framework"`
	EnvVariables map[string]string        `json:"envVars"`
	Capabilities CapabilityConfigurations `json:"capabilities"`
}

func convertCapabilities(parsed yaml.CapabilityConfigurations) CapabilityConfigurations {
	result := make(CapabilityConfigurations, len(parsed))
	for i, capValue := range parsed {
		moduleConstraint := "latest"
		if capValue.ModuleConstraint != nil {
			moduleConstraint = *capValue.ModuleConstraint
		}
		result[i] = &CapabilityConfiguration{
			Name:             capValue.Name,
			ModuleSource:     capValue.ModuleSource,
			ModuleConstraint: moduleConstraint,
			Variables:        convertVariables(capValue.Variables),
			Connections:      convertConnections(capValue.Connections),
			Namespace:        capValue.Namespace,
		}
	}
	return result
}

func convertAppConfigurations(parsed map[string]yaml.AppConfiguration) map[string]*AppConfiguration {
	apps := make(map[string]*AppConfiguration)
	for name, value := range parsed {
		bc := blockConfigFromYaml(name, value.BlockConfiguration, BlockTypeApplication, types.CategoryApp)
		apps[name] = &AppConfiguration{
			BlockConfiguration: *bc,
			Framework:          value.Framework,
			EnvVariables:       value.EnvVariables,
			Capabilities:       convertCapabilities(value.Capabilities),
		}
	}
	return apps
}

func (a *AppConfiguration) Initialize(ctx context.Context, resolver core.InitializeResolver, ic core.IacContext, pc core.ObjectPathContext) core.InitializeErrors {
	errs := a.BlockConfiguration.Initialize(ctx, resolver, ic, pc)
	errs = append(errs, a.Capabilities.Initialize(ctx, resolver, ic, pc, a.Module)...)
	return errs
}

func (a *AppConfiguration) Resolve(ctx context.Context, resolver core.ResolveResolver, ic core.IacContext, pc core.ObjectPathContext) core.ResolveErrors {
	errs := a.BlockConfiguration.Resolve(ctx, resolver, ic, pc)
	errs = append(errs, a.Capabilities.Resolve(ctx, resolver, ic, pc)...)
	return errs
}

func (a *AppConfiguration) Validate(ic core.IacContext, pc core.ObjectPathContext) core.ValidateErrors {
	errs := a.BlockConfiguration.Validate(ic, pc)
	errs = append(errs, a.ValidateEnvVariables(pc)...)
	errs = append(errs, a.Capabilities.Validate(ic, pc, a.Module)...)
	return errs
}

func hasInvalidChars(r rune) bool {
	return (r < 'A' || r > 'z') && r != '_' && (r < '0' || r > '9')
}

func startsWithNumber(s string) bool {
	return s[0] >= '0' && s[0] <= '9'
}

func (a *AppConfiguration) ValidateEnvVariables(pc core.ObjectPathContext) core.ValidateErrors {
	if len(a.EnvVariables) == 0 {
		return nil
	}

	errs := core.ValidateErrors{}
	for k, _ := range a.EnvVariables {
		curpc := pc.SubKey("environment", k)
		if startsWithNumber(k) {
			errs = append(errs, core.EnvVariableKeyStartsWithNumberError(curpc))
		}
		if strings.IndexFunc(k, hasInvalidChars) != -1 {
			errs = append(errs, core.EnvVariableKeyInvalidCharsError(curpc))
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (a *AppConfiguration) Normalize(ctx context.Context, pc core.ObjectPathContext, resolver core.ConnectionResolver) core.NormalizeErrors {
	errs := core.NormalizeErrors{}
	errs = append(errs, a.Connections.Normalize(ctx, pc, resolver)...)
	errs = append(errs, a.Capabilities.Normalize(ctx, pc, resolver)...)
	return errs
}

func (a *AppConfiguration) ToBlock(orgName string, stackId int64) types.Block {
	block := a.BlockConfiguration.ToBlock(orgName, stackId)
	block.Framework = a.Framework
	return block
}

func (a *AppConfiguration) ApplyChangesTo(ic core.IacContext, updater core.WorkspaceConfigUpdater) error {
	if err := a.BlockConfiguration.ApplyChangesTo(ic, updater); err != nil {
		return err
	}

	if ic.IsOverrides {
		for name, value := range a.EnvVariables {
			updater.AddOrUpdateEnvVariable(name, value, false)
		}
	} else {
		updater.RemoveEnvVariablesNotIn(a.EnvVariables)
		for name, value := range a.EnvVariables {
			updater.AddOrUpdateEnvVariable(name, value, false)
		}
	}

	return a.Capabilities.ApplyChangesTo(ic, updater)
}
