package config

import (
	"context"
	"github.com/BSick7/go-api/errors"
	"github.com/nullstone-io/iac/core"
	"github.com/nullstone-io/iac/yaml"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
	"strings"
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

func convertAppConfigurations(parsed map[string]yaml.AppConfiguration) map[string]*AppConfiguration {
	apps := make(map[string]*AppConfiguration)
	for name, value := range parsed {
		bc := blockConfigFromYaml(name, value.BlockConfiguration, BlockTypeApplication, types.CategoryApp)
		apps[name] = &AppConfiguration{
			BlockConfiguration: *bc,
			EnvVariables:       value.EnvVariables,
			Capabilities:       convertCapabilities(value.Capabilities),
		}
	}
	return apps
}

func (a *AppConfiguration) Resolve(ctx context.Context, resolver core.ModuleVersionResolver, ic core.IacContext, pc core.ObjectPathContext) core.ResolveErrors {
	errs := a.BlockConfiguration.Resolve(ctx, resolver, ic, pc)
	errs = append(errs, a.ResolveCapabilities(ctx, resolver, ic, pc)...)
	return errs
}

func (a *AppConfiguration) ResolveCapabilities(ctx context.Context, resolver core.ModuleVersionResolver, ic core.IacContext, pc core.ObjectPathContext) core.ResolveErrors {
	if len(a.Capabilities) == 0 {
		return nil
	}
	errs := core.ResolveErrors{}
	for i, iacCap := range a.Capabilities {
		curpc := pc.SubIndex("capabilities", i)
		var err *core.ResolveError
		if a.Capabilities[i], err = a.ResolveCapability(ctx, resolver, ic, curpc, iacCap); err != nil {
			errs = append(errs, *err)
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (a *AppConfiguration) ResolveCapability(ctx context.Context, resolver core.ModuleVersionResolver, ic core.IacContext, pc core.ObjectPathContext, iacCap CapabilityConfiguration) (CapabilityConfiguration, *core.ResolveError) {
	if ic.IsOverrides && iacCap.ModuleSource == "" {
		// TODO: Add support for loading module in overrides file
		return iacCap, nil
	}

	contract := types.ModuleContractName{
		Category: string(types.CategoryCapability),
		Provider: "*",
		Platform: "*",
	}
	if a.Module != nil {
		contract.Provider = strings.Join(a.Module.ProviderTypes, ",")
	}
	m, mv, err := core.ResolveModule(ctx, resolver, pc, iacCap.ModuleSource, iacCap.ModuleSourceVersion, contract)
	if err != nil {
		return iacCap, err
	}
	iacCap.Module = m
	iacCap.ModuleVersion = mv
	return iacCap, nil
}

func (a *AppConfiguration) Validate(ctx context.Context, resolver core.ValidateResolver, ic core.IacContext, pc core.ObjectPathContext) errors.ValidationErrors {
	ve := a.BlockConfiguration.Validate(ctx, resolver, ic, pc)
	ve = append(ve, a.ValidateEnvVariables(ic, pc)...)
	ve = append(ve, a.Capabilities.Validate(ctx, resolver, ic, pc, a.Module)...)
	return ve
}

func hasInvalidChars(r rune) bool {
	return (r < 'A' || r > 'z') && r != '_' && (r < '0' || r > '9')
}

func startsWithNumber(s string) bool {
	return s[0] >= '0' && s[0] <= '9'
}

func (a *AppConfiguration) ValidateEnvVariables(ic core.IacContext, pc core.ObjectPathContext) errors.ValidationErrors {
	if len(a.EnvVariables) == 0 {
		return nil
	}

	ve := errors.ValidationErrors{}
	for k, _ := range a.EnvVariables {
		curpc := pc.SubKey("environment", k)
		if startsWithNumber(k) {
			ve = append(ve, EnvVariableKeyStartsWithNumberError(ic, curpc))
		}
		if strings.IndexFunc(k, hasInvalidChars) != -1 {
			ve = append(ve, EnvVariableKeyInvalidCharsError(ic, curpc))
		}
	}

	if len(ve) > 0 {
		return ve
	}
	return nil
}

func (a *AppConfiguration) Normalize(ctx context.Context, resolver core.ConnectionResolver) error {
	err := NormalizeConnectionTargets(ctx, a.Connections, resolver)
	if err != nil {
		return err
	}
	return a.Capabilities.Normalize(ctx, resolver)
}

func (a *AppConfiguration) ToBlock(orgName string, stackId int64) types.Block {
	block := a.BlockConfiguration.ToBlock(orgName, stackId)
	block.Capabilities = a.Capabilities.ToCapabilities(stackId)
	return block
}
