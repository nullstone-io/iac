package config

import (
	"context"
	"github.com/nullstone-io/iac/core"
	"github.com/nullstone-io/module/config"
	"gopkg.in/nullstone-io/go-api-client.v0/artifacts"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
	"strings"
)

type ConnectionConfigurations map[string]*ConnectionConfiguration

func (s ConnectionConfigurations) DesiredTargets() types.ConnectionTargets {
	targets := types.ConnectionTargets{}
	for key, c := range s {
		targets[key] = c.DesiredTarget
	}
	return targets
}

func (s ConnectionConfigurations) Resolve(ctx context.Context, resolver core.ResolveResolver, ic core.IacContext, pc core.ObjectPathContext,
	blockManifest config.Manifest) core.ResolveErrors {
	var errs core.ResolveErrors
	if !ic.IsOverrides {
		for name, manifestConn := range blockManifest.Connections {
			if _, inNsConfig := s[name]; !inNsConfig && !manifestConn.Optional {
				errs = append(errs, core.MissingRequiredConnectionError(pc.SubField("connections"), name))
			}
		}
	}
	for name, c := range s {
		if schema, ok := blockManifest.Connections[name]; ok {
			c.Schema = &schema
		}
		if err := c.Resolve(ctx, resolver, pc.SubKey("connections", name)); err != nil {
			errs = append(errs, *err)
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

// Validate performs validation on all IaC connections by matching them against connections in the module
func (s ConnectionConfigurations) Validate(pc core.ObjectPathContext, moduleName string) core.ValidateErrors {
	if len(s) == 0 {
		return nil
	}
	errs := core.ValidateErrors{}
	for name, cb := range s {
		if err := cb.Validate(pc.SubKey("connections", name), moduleName); err != nil {
			errs = append(errs, *err)
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

// Normalize loops through all connections and does the following to DesiredTarget:
// 1. Fills all fields (Id+Name for Stack/Block/Env)
// 2. If block.IsShared, resolves the Env to the previews-shared env
func (s ConnectionConfigurations) Normalize(ctx context.Context, pc core.ObjectPathContext, resolver core.ConnectionResolver) core.NormalizeErrors {
	errs := core.NormalizeErrors{}
	for name, connection := range s {
		if err := connection.Normalize(ctx, pc.SubKey("connections", name), resolver); err != nil {
			errs = append(errs, *err)
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

type ConnectionConfiguration struct {
	DesiredTarget   types.ConnectionTarget `json:"desiredTarget"`
	EffectiveTarget types.ConnectionTarget `json:"effectiveTarget"`
	Schema          *config.Connection     `json:"schema"`
	Block           *types.Block           `json:"block"`
	Module          *types.Module          `json:"module"`
}

// Resolve resolves the connection's target (i.e., block) and matches the connection contract
func (c *ConnectionConfiguration) Resolve(ctx context.Context, resolver core.ResolveResolver, pc core.ObjectPathContext) *core.ResolveError {
	if c.Schema == nil || c.DesiredTarget.BlockName == "" {
		// There is nothing to resolve
		// Validate will report errors
		return nil
	}

	found, err := resolver.ResolveBlock(ctx, c.DesiredTarget)
	if err != nil {
		if find.IsMissingResource(err) {
			return core.MissingConnectionTargetError(pc, err)
		}
		return core.LookupConnectionTargetFailedError(pc, err)
	}
	c.Block = &found

	if found.ModuleSource == "" {
		return core.ResolvedBlockMissingModuleError(pc, c.DesiredTarget.StackName, c.DesiredTarget.BlockName)
	}

	ms, err := artifacts.ParseSource(found.ModuleSource)
	if err != nil {
		return core.InvalidModuleFormatError(pc, found.ModuleSource)
	}
	m, mErr := resolver.ResolveModule(ctx, *ms)
	if mErr != nil {
		return core.ModuleLookupFailedError(pc, found.ModuleSource, mErr)
	}
	c.Module = m
	return nil
}

func (c *ConnectionConfiguration) Validate(pc core.ObjectPathContext, moduleName string) *core.ValidateError {
	if c.Schema == nil {
		return core.ConnectionDoesNotExistError(pc, moduleName)
	}
	if c.Schema.Optional && c.DesiredTarget.IsEmpty() {
		return nil
	}
	if c.DesiredTarget.BlockName == "" {
		return core.MissingConnectionBlockError(pc)
	}
	if c.Module == nil {
		return nil
	}

	mcn1, mcnErr := types.ParseModuleContractName(c.Schema.Contract)
	if mcnErr != nil {
		return core.InvalidConnectionContractError(pc, c.Schema.Contract, moduleName)
	}
	mcn2 := types.ModuleContractName{
		Category:    string(c.Module.Category),
		Subcategory: string(c.Module.Subcategory),
		Provider:    strings.Join(c.Module.ProviderTypes, ","),
		Platform:    c.Module.Platform,
		Subplatform: c.Module.Subplatform,
	}
	if ok := mcn1.Match(mcn2); !ok {
		return core.MismatchedConnectionContractError(pc, c.Block.Name, c.Schema.Contract)
	}

	return nil
}

func (c *ConnectionConfiguration) Normalize(ctx context.Context, pc core.ObjectPathContext, resolver core.ConnectionResolver) *core.NormalizeError {
	if c.DesiredTarget.IsEmpty() {
		c.EffectiveTarget = types.ConnectionTarget{}
		return nil
	}
	ct, err := resolver.ResolveConnection(ctx, c.DesiredTarget)
	if err != nil {
		return &core.NormalizeError{
			ObjectPathContext: pc,
			ErrorMessage:      err.Error(),
		}
	}
	c.DesiredTarget.StackId = ct.StackId
	c.DesiredTarget.StackName = ct.StackName
	c.DesiredTarget.BlockId = ct.BlockId
	c.DesiredTarget.BlockName = ct.BlockName
	if c.DesiredTarget.EnvId != nil {
		c.DesiredTarget.EnvName = ct.EnvName
	} else if c.DesiredTarget.EnvName != "" {
		c.DesiredTarget.EnvId = ct.EnvId
	}
	c.EffectiveTarget = ct
	return nil
}
