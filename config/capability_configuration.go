package config

import (
	"context"
	"fmt"
	"github.com/nullstone-io/iac/core"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

var (
	_ core.ChangeApplier = CapabilityConfigurations{}
	_ core.ChangeApplier = CapabilityConfiguration{}
)

type CapabilityConfigurations []CapabilityConfiguration

func (c CapabilityConfigurations) Identities() []core.CapabilityIdentity {
	result := make([]core.CapabilityIdentity, 0)
	for _, cur := range c {
		result = append(result, cur.Identity())
	}
	return result
}

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
	pc core.ObjectPathContext, appModule *types.Module) core.ValidateErrors {
	if len(c) == 0 {
		return nil
	}
	errs := core.ValidateErrors{}
	for i, iacCap := range c {
		curpc := pc.SubIndex("capabilities", i)
		errs = append(errs, iacCap.Validate(ctx, resolver, ic, curpc, appModule)...)
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (c CapabilityConfigurations) ToCapabilities(stackId int64) []types.Capability {
	var result []types.Capability
	for _, cur := range c {
		capability := types.Capability{
			ModuleSource:        cur.ModuleSource,
			ModuleSourceVersion: cur.ModuleSourceVersion,
			Connections:         types.ConnectionTargets{},
		}
		if cur.Namespace != nil {
			capability.Namespace = *cur.Namespace
		}
		for key, conn := range cur.Connections {
			target := conn
			if target.StackId == 0 && target.StackName == "" {
				target.StackId = stackId
			}
			capability.Connections[key] = target
		}
		result = append(result, capability)
	}
	return result
}

func (c CapabilityConfigurations) ApplyChangesTo(ic core.IacContext, updater core.WorkspaceConfigUpdater) error {
	if ic.IsOverrides {
		for _, cur := range c {
			if err := cur.ApplyChangesTo(ic, updater); err != nil {
				return err
			}
		}
	} else {
		updater.RemoveCapabilitiesNotIn(c.Identities())
		for _, cur := range c {
			if err := cur.ApplyChangesTo(ic, updater); err != nil {
				return err
			}
		}
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

func (c CapabilityConfiguration) Identity() core.CapabilityIdentity {
	return core.CapabilityIdentity{
		ModuleSource:      c.ModuleSource,
		ConnectionTargets: c.Connections,
	}
}

func (c CapabilityConfiguration) Normalize(ctx context.Context, resolver core.ConnectionResolver) (CapabilityConfiguration, error) {
	if err := NormalizeConnectionTargets(ctx, c.Connections, resolver); err != nil {
		return c, err
	}
	return c, nil
}

func (c CapabilityConfiguration) Validate(ctx context.Context, resolver core.ValidateResolver, ic core.IacContext,
	pc core.ObjectPathContext, appModule *types.Module) core.ValidateErrors {
	if c.Module == nil {
		// We can't perform validation if the module isn't loaded
		return nil
	}
	if ic.IsOverrides && c.ModuleSource == "" {
		// TODO: Add support for validating variables and connections in an overrides file
		return nil
	}

	errs := core.ValidateErrors{}
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
			errs = append(errs, UnsupportedAppCategoryError(pc.SubField("module"), c.ModuleSource, string(appModule.Subcategory)))
		}
	}

	// if we were able to find the module version
	//   1. validate each of the variables to ensure the module supports them
	//   2. validate each of the connections to ensure the block matches the connection contract
	if mv := c.ModuleVersion; mv != nil {
		moduleName := fmt.Sprintf("%s@%s", c.ModuleSource, c.ModuleSourceVersion)
		errs = append(errs, core.ValidateVariables(pc, c.Variables, mv.Manifest.Variables, moduleName)...)
		errs = append(errs, core.ValidateConnections(ctx, resolver, pc, c.Connections, mv.Manifest.Connections, moduleName)...)
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (c CapabilityConfiguration) ApplyChangesTo(ic core.IacContext, updater core.WorkspaceConfigUpdater) error {
	capUpdater := updater.GetCapabilityUpdater(c.Identity())
	if capUpdater == nil {
		return nil
	}

	capUpdater.UpdateSchema(c.ModuleSource, c.ModuleVersion)
	capUpdater.UpdateNamespace(c.Namespace)
	for name, value := range c.Variables {
		capUpdater.UpdateVariableValue(name, value)
	}
	for name, value := range c.Connections {
		capUpdater.UpdateConnectionTarget(name, value)
	}
	return nil
}
