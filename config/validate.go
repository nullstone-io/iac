package config

import (
	"context"
	"github.com/BSick7/go-api/errors"
	"github.com/nullstone-io/iac/core"
	"github.com/nullstone-io/module/config"
	"gopkg.in/nullstone-io/go-api-client.v0/artifacts"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
	"strings"
)

func ValidateVariables(ic core.IacContext, pc core.YamlPathContext, variables map[string]any, expectedVariables map[string]config.Variable, moduleName string) errors.ValidationErrors {
	if len(variables) == 0 {
		return nil
	}
	ve := errors.ValidationErrors{}
	for k, _ := range variables {
		if _, ok := expectedVariables[k]; !ok {
			ve = append(ve, VariableDoesNotExistError(ic, pc.SubKey("vars", k), moduleName))
		}
	}
	if len(ve) > 0 {
		return ve
	}
	return nil
}

// ValidateConnections performs validation on all IaC connections by matching them against connections in the module
func ValidateConnections(ctx context.Context, resolver core.ValidateResolver, ic core.IacContext, pc core.YamlPathContext,
	connections types.ConnectionTargets, expectedConnections map[string]config.Connection, moduleName string) errors.ValidationErrors {
	if len(connections) == 0 {
		return nil
	}
	ve := errors.ValidationErrors{}
	for key, conn := range connections {
		curpc := pc.SubKey("connections", key)
		manifestConnection, found := expectedConnections[key]
		if !found {
			ve = append(ve, ConnectionDoesNotExistError(ic, curpc, moduleName))
			continue
		}
		verr := ValidateConnection(ctx, resolver, ic, curpc, conn, manifestConnection, moduleName)
		if verr != nil {
			ve = append(ve, *verr)
		}
	}
	if len(ve) > 0 {
		return ve
	}
	return nil
}

// ValidateConnection performs validation on a single IaC connection after it has been matched to a connection in the module manifest
//  1. Verifies that a connection specified in IaC exists in the module
//  2. Resolves the connection's target (i.e. block)
//  3. Verifies the block matches the connection contract
func ValidateConnection(ctx context.Context, resolver core.ValidateResolver, ic core.IacContext, pc core.YamlPathContext, connection types.ConnectionTarget, manifestConnection config.Connection, moduleName string) *errors.ValidationError {
	if connection.BlockName == "" {
		return MissingConnectionBlockError(ic, pc)
	}

	found, err := resolver.ResolveBlock(ctx, connection)
	if err != nil {
		if find.IsMissingResource(err) {
			return MissingConnectionTargetError(ic, pc, err)
		}
		return LookupConnectionTargetFailedError(ic, pc, err)
	}

	mcn1, mcnErr := types.ParseModuleContractName(manifestConnection.Contract)
	if mcnErr != nil {
		return InvalidConnectionContractError(ic, pc, manifestConnection.Contract, moduleName)
	}
	ms, err := artifacts.ParseSource(found.ModuleSource)
	if err != nil {
		return InvalidModuleFormatError(ic, pc, found.ModuleSource)
	}
	m, mErr := resolver.ResolveModule(ctx, *ms)
	if mErr != nil {
		return ModuleLookupFailedError(ic, pc, found.ModuleSource, mErr)
	}
	if mcnErr == nil && m != nil {
		mcn2 := types.ModuleContractName{
			Category:    string(m.Category),
			Subcategory: string(m.Subcategory),
			Provider:    strings.Join(m.ProviderTypes, ","),
			Platform:    m.Platform,
			Subplatform: m.Subplatform,
		}
		if ok := mcn1.Match(mcn2); !ok {
			return MismatchedConnectionContractError(ic, pc, found.Name, manifestConnection.Contract)
		}
	}

	return nil
}
