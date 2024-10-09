package config

import (
	"context"
	"fmt"
	"github.com/BSick7/go-api/errors"
	"github.com/nullstone-io/iac/core"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

// ValidateCapabilities performs validation on all IaC capabilities within an application
func ValidateCapabilities(ctx context.Context, resolver *find.ResourceResolver, ic core.IacContext, pc core.YamlPathContext, capabilities CapabilityConfigurations, subcategory types.SubcategoryName) errors.ValidationErrors {
	if len(capabilities) == 0 {
		return nil
	}
	ve := errors.ValidationErrors{}
	for i, iacCap := range capabilities {
		curpc := pc.SubIndex("capabilities", i)
		ve = append(ve, ValidateCapability(ctx, resolver, ic, curpc, iacCap, string(subcategory))...)
	}

	if len(ve) > 0 {
		return ve
	}
	return nil
}

func ValidateCapability(ctx context.Context, resolver *find.ResourceResolver, ic core.IacContext, pc core.YamlPathContext, iacCap CapabilityConfiguration, subcategory string) errors.ValidationErrors {
	// ensure the module is a capability module and supports the provider type (e.g. aws, gcp, azure)
	providerType, err := resolver.ResolveCurProviderType(ctx)
	if err != nil {
		return errors.ValidationErrors{LookupProviderTypeFailedError(ic, pc, err)}
	}
	contract := fmt.Sprintf("capability/%s/*", providerType)
	m, mv, verr := ResolveModule(ctx, resolver, ic, pc, iacCap.ModuleSource, iacCap.ModuleSourceVersion, contract)
	if verr != nil {
		return errors.ValidationErrors{*verr}
	}

	ve := errors.ValidationErrors{}
	// check to make sure the capability module supports the subcategory
	// examples are "container", "serverless", "static-site", "server"
	if m != nil {
		found := false
		for _, cat := range m.AppCategories {
			if cat == subcategory {
				found = true
				break
			}
		}
		if !found {
			ve = append(ve, UnsupportedAppCategoryError(ic, pc.SubField("module"), iacCap.ModuleSource, subcategory))
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
