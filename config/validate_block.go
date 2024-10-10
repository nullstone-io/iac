package config

import (
	"context"
	"fmt"
	"github.com/BSick7/go-api/errors"
	"github.com/nullstone-io/iac/core"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

func ValidateBlock(ctx context.Context, resolver *find.ResourceResolver, ic core.IacContext, ipc core.YamlPathContext,
	contract, moduleSource, moduleSourceVersion string, variables map[string]any, connections types.ConnectionTargets,
	envVars map[string]string, capabilities CapabilityConfigurations) errors.ValidationErrors {

	ve := errors.ValidationErrors{}

	var subcategory types.SubcategoryName
	skipModuleCheck := ic.IsOverrides && moduleSource == ""
	if !skipModuleCheck {
		// TODO: Add support for validating variables and connections in an overrides file
		m, mv, err := ResolveModule(ctx, resolver, ic, ipc, moduleSource, moduleSourceVersion, contract)
		if err != nil {
			return errors.ValidationErrors{*err}
		}
		subcategory = m.Subcategory

		moduleName := fmt.Sprintf("%s/%s@%s", m.OrgName, m.Name, mv.Version)
		ve = append(ve, ValidateVariables(ic, ipc, variables, mv.Manifest.Variables, moduleName)...)
		ve = append(ve, ValidateConnections(ctx, resolver, ic, ipc, connections, mv.Manifest.Connections, moduleName)...)
	}

	ve = append(ve, ValidateEnvVariables(ic, ipc, envVars)...)
	ve = append(ve, ValidateCapabilities(ctx, resolver, ic, ipc, capabilities, subcategory)...)

	if len(ve) > 0 {
		return ve
	}
	return nil
}