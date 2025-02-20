package config

import (
	"context"
	"fmt"
	"github.com/nullstone-io/iac/core"
	"github.com/nullstone-io/module/config"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
	"strings"
)

var (
	_ core.ChangeApplier = CapabilityConfigurations{}
	_ core.ChangeApplier = &CapabilityConfiguration{}
)

type CapabilityConfigurations []*CapabilityConfiguration

func (c CapabilityConfigurations) Identities() []core.CapabilityIdentity {
	result := make([]core.CapabilityIdentity, 0)
	for _, cur := range c {
		result = append(result, cur.Identity())
	}
	return result
}

func (c CapabilityConfigurations) Normalize(ctx context.Context, pc core.ObjectPathContext, resolver core.ConnectionResolver) core.NormalizeErrors {
	for i, iacCap := range c {
		if err := iacCap.Connections.Normalize(ctx, pc.SubIndex("capabilities", i), resolver); err != nil {
			return err
		}
	}
	return nil
}

// Validate performs validation on all IaC capabilities
func (c CapabilityConfigurations) Validate(ic core.IacContext, pc core.ObjectPathContext, appModule *types.Module) core.ValidateErrors {
	if len(c) == 0 {
		return nil
	}
	errs := core.ValidateErrors{}
	for i, iacCap := range c {
		errs = append(errs, iacCap.Validate(ic, pc.SubIndex("capabilities", i), appModule)...)
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (c CapabilityConfigurations) ToCapabilities() []types.Capability {
	var result []types.Capability
	for _, cur := range c {
		capability := types.Capability{
			Name:                cur.Name,
			ModuleSource:        cur.ModuleSource,
			ModuleSourceVersion: cur.ModuleConstraint,
			Connections:         cur.Connections.Targets(),
		}
		if cur.Namespace != nil {
			capability.Namespace = *cur.Namespace
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

func (c CapabilityConfigurations) Resolve(ctx context.Context, resolver core.ResolveResolver, ic core.IacContext,
	pc core.ObjectPathContext, appModule *types.Module) core.ResolveErrors {
	if len(c) == 0 {
		return nil
	}
	errs := core.ResolveErrors{}
	for i, iacCap := range c {
		errs = append(errs, iacCap.Resolve(ctx, resolver, ic, pc.SubIndex("capabilities", i), appModule)...)
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

type CapabilityConfiguration struct {
	Name             string                   `json:"name"`
	ModuleSource     string                   `json:"moduleSource"`
	ModuleConstraint string                   `json:"moduleConstraint"`
	Variables        VariableConfigurations   `json:"vars"`
	Connections      ConnectionConfigurations `json:"connections"`
	Namespace        *string                  `json:"namespace"`

	Module        *types.Module        `json:"module"`
	ModuleVersion *types.ModuleVersion `json:"moduleVersion"`
}

func (c *CapabilityConfiguration) Identity() core.CapabilityIdentity {
	return core.CapabilityIdentity{
		Name:              c.Name,
		ModuleSource:      c.ModuleSource,
		ConnectionTargets: c.Connections.Targets(),
	}
}

func (c *CapabilityConfiguration) Resolve(ctx context.Context, resolver core.ResolveResolver, ic core.IacContext,
	pc core.ObjectPathContext, appModule *types.Module) core.ResolveErrors {
	if c.Variables == nil {
		c.Variables = VariableConfigurations{}
	}
	if c.Connections == nil {
		c.Connections = ConnectionConfigurations{}
	}
	if ic.IsOverrides && c.ModuleSource == "" {
		// TODO: Add support for loading module in overrides file
		return nil
	}

	errs := core.ResolveErrors{}

	contract := types.ModuleContractName{
		Category: string(types.CategoryCapability),
		Provider: "*",
		Platform: "*",
	}
	if appModule != nil {
		contract.Provider = strings.Join(appModule.ProviderTypes, ",")
	}

	manifest := config.Manifest{Variables: map[string]config.Variable{}, Connections: map[string]config.Connection{}}
	m, mv, err := core.ResolveModule(ctx, resolver, pc, c.ModuleSource, c.ModuleConstraint, contract)
	if err != nil {
		errs = append(errs, *err)
	} else {
		c.Module = m
		c.ModuleVersion = mv
		manifest = mv.Manifest
	}
	errs = append(errs, c.Variables.Resolve(manifest)...)
	errs = append(errs, c.Connections.Resolve(ctx, resolver, pc, manifest)...)
	return errs
}

func (c *CapabilityConfiguration) Validate(ic core.IacContext, pc core.ObjectPathContext, appModule *types.Module) core.ValidateErrors {
	errs := core.ValidateErrors{}
	// TODO: After deprecating (ModuleSource+ConnectionTargets), validate Name
	//if c.Name == "" {
	//	err := core.MissingCapabilityNameError(pc)
	//	errs = append(errs, *err)
	//}

	if c.Module == nil {
		// We can't perform validation if the module isn't loaded
		return errs
	}
	if ic.IsOverrides && c.ModuleSource == "" {
		// TODO: Add support for validating variables and connections in an overrides file
		return errs
	}

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
			errs = append(errs, core.UnsupportedAppCategoryError(pc.SubField("module"), c.ModuleSource, string(appModule.Subcategory)))
		}
	}

	//   1. validate each of the variables to ensure the module supports them
	//   2. validate each of the connections to ensure the block matches the connection contract
	moduleName := fmt.Sprintf("%s@%s", c.ModuleSource, c.ModuleConstraint)
	errs = append(errs, c.Variables.Validate(pc, moduleName)...)
	errs = append(errs, c.Connections.Validate(pc, moduleName)...)
	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (c *CapabilityConfiguration) ApplyChangesTo(ic core.IacContext, updater core.WorkspaceConfigUpdater) error {
	capUpdater := updater.GetCapabilityUpdater(c.Identity())
	if capUpdater == nil {
		return nil
	}

	capUpdater.UpdateSchema(c.ModuleSource, c.ModuleVersion)
	capUpdater.UpdateNamespace(c.Namespace)
	for name, vc := range c.Variables {
		capUpdater.UpdateVariableValue(name, vc.Value)
	}
	for name, cc := range c.Connections {
		capUpdater.UpdateConnectionTarget(name, cc.DesiredTarget, cc.EffectiveTarget)
	}
	return nil
}
