package config

import (
	"context"
	"fmt"
	"github.com/BSick7/go-api/errors"
	"github.com/nullstone-io/iac/core"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

type CapabilityConfigurations []CapabilityConfiguration

func (c CapabilityConfigurations) Normalize(ctx context.Context, resolver core.ConnectionResolver) error {
	for i, iacCap := range c {
		resolved, err := iacCap.Normalize(ctx, resolver)
		if err != nil {
			return err
		}
		c[i] = resolved
	}
	return nil
}

// Validate performs validation on all IaC capabilities
func (c CapabilityConfigurations) Validate(ctx context.Context, resolver core.ValidateResolver, ic core.IacContext,
	pc core.ObjectPathContext, appModule *types.Module) errors.ValidationErrors {
	if len(c) == 0 {
		return nil
	}
	ve := errors.ValidationErrors{}
	for i, iacCap := range c {
		curpc := pc.SubIndex("capabilities", i)
		ve = append(ve, iacCap.Validate(ctx, resolver, ic, curpc, appModule)...)
	}

	if len(ve) > 0 {
		return ve
	}
	return nil
}

type CapabilityConfiguration struct {
	ModuleSource        string                  `json:"moduleSource"`
	ModuleSourceVersion string                  `json:"moduleSourceVersion"`
	Variables           map[string]any          `json:"vars"`
	Connections         types.ConnectionTargets `json:"connections"`
	Namespace           *string                 `json:"namespace"`

	Module        *types.Module        `json:"module"`
	ModuleVersion *types.ModuleVersion `json:"moduleVersion"`
}

func (c CapabilityConfiguration) Normalize(ctx context.Context, resolver core.ConnectionResolver) (CapabilityConfiguration, error) {
	if err := NormalizeConnectionTargets(ctx, c.Connections, resolver); err != nil {
		return c, err
	}
	return c, nil
}

func (c CapabilityConfiguration) Validate(ctx context.Context, resolver core.ValidateResolver, ic core.IacContext,
	pc core.ObjectPathContext, appModule *types.Module) errors.ValidationErrors {
	if c.Module == nil {
		// We can't perform validation if the module isn't loaded
		return nil
	}
	if ic.IsOverrides && c.ModuleSource == "" {
		// TODO: Add support for validating variables and connections in an overrides file
		return nil
	}

	ve := errors.ValidationErrors{}
	// check to make sure the capability module supports the subcategory
	// examples are "container", "serverless", "static-site", "server"
	// TODO: Add support for validating app category
	if appModule != nil {
		found := false
		for _, cat := range c.Module.AppCategories {
			if cat == string(appModule.Subcategory) {
				found = true
				break
			}
		}
		if !found {
			ve = append(ve, UnsupportedAppCategoryError(ic, pc.SubField("module"), c.ModuleSource, string(appModule.Subcategory)))
		}
	}

	// if we were able to find the module version
	//   1. validate each of the variables to ensure the module supports them
	//   2. validate each of the connections to ensure the block matches the connection contract
	if mv := c.ModuleVersion; mv != nil {
		moduleName := fmt.Sprintf("%s@%s", c.ModuleSource, c.ModuleSourceVersion)
		ve = append(ve, ValidateVariables(ic, pc, c.Variables, mv.Manifest.Variables, moduleName)...)
		ve = append(ve, ValidateConnections(ctx, resolver, ic, pc, c.Connections, mv.Manifest.Connections, moduleName)...)
	}
	if len(ve) > 0 {
		return ve
	}
	return nil
}
