package config

import (
	"context"
	"fmt"
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

func convertAppConfigurations(parsed map[string]yaml.AppConfiguration) map[string]AppConfiguration {
	apps := make(map[string]AppConfiguration)
	for appName, appValue := range parsed {
		app := AppConfiguration{
			BlockConfiguration: blockConfigFromYaml(appName, appValue.BlockConfiguration, BlockTypeApplication, types.CategoryApp),
			EnvVariables:       appValue.EnvVariables,
			Capabilities:       convertCapabilities(appValue.Capabilities),
		}
		apps[appName] = app
	}
	return apps
}

func (a *AppConfiguration) Validate(ctx context.Context, resolver core.ValidateResolver, ic core.IacContext, pc core.YamlPathContext) errors.ValidationErrors {
	ve := a.BlockConfiguration.Validate(ctx, resolver, ic, pc)
	ve = append(ve, a.ValidateEnvVariables(ic, pc)...)
	ve = append(ve, a.ValidateCapabilities(ctx, resolver, ic, pc)...)
	return ve
}

// ValidateCapabilities performs validation on all IaC capabilities within an application
func (a *AppConfiguration) ValidateCapabilities(ctx context.Context, resolver core.ValidateResolver, ic core.IacContext, pc core.YamlPathContext) errors.ValidationErrors {
	if len(a.Capabilities) == 0 {
		return nil
	}
	ve := errors.ValidationErrors{}
	for i, iacCap := range a.Capabilities {
		curpc := pc.SubIndex("capabilities", i)
		ve = append(ve, a.ValidateCapability(ctx, resolver, ic, curpc, iacCap)...)
	}

	if len(ve) > 0 {
		return ve
	}
	return nil
}

func (a *AppConfiguration) ValidateCapability(ctx context.Context, resolver core.ValidateResolver, ic core.IacContext, pc core.YamlPathContext, iacCap CapabilityConfiguration) errors.ValidationErrors {
	if ic.IsOverrides && iacCap.ModuleSource == "" {
		// TODO: Add support for validating variables and connections in an overrides file
		return nil
	}
	if a.Module == nil {
		// The app module isn't loaded, we can't perform validation on this capability
		return nil
	}

	contract := types.ModuleContractName{
		Category: string(types.CategoryCapability),
		Provider: strings.Join(a.Module.ProviderTypes, ","),
		Platform: "*",
	}
	m, mv, verr := ResolveModule(ctx, resolver, ic, pc, iacCap.ModuleSource, iacCap.ModuleSourceVersion, contract)
	if verr != nil {
		return errors.ValidationErrors{*verr}
	}

	ve := errors.ValidationErrors{}
	// check to make sure the capability module supports the subcategory
	// examples are "container", "serverless", "static-site", "server"
	skipAppCategoryCheck := ic.IsOverrides && a.Module.Subcategory == ""
	// TODO: Add support for validating app category
	if m != nil && !skipAppCategoryCheck {
		found := false
		for _, cat := range m.AppCategories {
			if cat == string(a.Module.Subcategory) {
				found = true
				break
			}
		}
		if !found {
			ve = append(ve, UnsupportedAppCategoryError(ic, pc.SubField("module"), iacCap.ModuleSource, string(a.Module.Subcategory)))
		}
	}

	// if we were able to find the module version
	//   1. validate each of the variables to ensure the module supports them
	//   2. validate each of the connections to ensure the block matches the connection contract
	if mv != nil {
		moduleName := fmt.Sprintf("%s@%s", iacCap.ModuleSource, iacCap.ModuleSourceVersion)
		ve = append(ve, ValidateVariables(ic, pc, iacCap.Variables, mv.Manifest.Variables, moduleName)...)
		ve = append(ve, ValidateConnections(ctx, resolver, ic, pc, iacCap.Connections, mv.Manifest.Connections, moduleName)...)
	}

	if len(ve) > 0 {
		return ve
	}
	return nil
}

func hasInvalidChars(r rune) bool {
	return (r < 'A' || r > 'z') && r != '_' && (r < '0' || r > '9')
}

func startsWithNumber(s string) bool {
	return s[0] >= '0' && s[0] <= '9'
}

func (a *AppConfiguration) ValidateEnvVariables(ic core.IacContext, pc core.YamlPathContext) errors.ValidationErrors {
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
