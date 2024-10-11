package core

import (
	"context"
	"github.com/nullstone-io/module/config"
	"gopkg.in/nullstone-io/go-api-client.v0/artifacts"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
	"strings"
)

func ValidateVariables(pc ObjectPathContext, variables map[string]any, expectedVariables map[string]config.Variable, moduleName string) ValidateErrors {
	if len(variables) == 0 {
		return nil
	}
	ve := ValidateErrors{}
	for k, _ := range variables {
		if _, ok := expectedVariables[k]; !ok {
			ve = append(ve, VariableDoesNotExistError(pc.SubKey("vars", k), moduleName))
		}
	}
	if len(ve) > 0 {
		return ve
	}
	return nil
}

// ValidateConnections performs validation on all IaC connections by matching them against connections in the module
func ValidateConnections(ctx context.Context, resolver ValidateResolver, pc ObjectPathContext,
	connections types.ConnectionTargets, expectedConnections map[string]config.Connection, moduleName string) ValidateErrors {
	if len(connections) == 0 {
		return nil
	}
	ve := ValidateErrors{}
	for key, conn := range connections {
		curpc := pc.SubKey("connections", key)
		manifestConnection, found := expectedConnections[key]
		if !found {
			ve = append(ve, ConnectionDoesNotExistError(curpc, moduleName))
			continue
		}
		verr := ValidateConnection(ctx, resolver, curpc, conn, manifestConnection, moduleName)
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
func ValidateConnection(ctx context.Context, resolver ValidateResolver, pc ObjectPathContext, connection types.ConnectionTarget, manifestConnection config.Connection, moduleName string) *ValidateError {
	if connection.BlockName == "" {
		return MissingConnectionBlockError(pc)
	}

	found, err := resolver.ResolveBlock(ctx, connection)
	if err != nil {
		if find.IsMissingResource(err) {
			return MissingConnectionTargetError(pc, err)
		}
		return LookupConnectionTargetFailedError(pc, err)
	}

	mcn1, mcnErr := types.ParseModuleContractName(manifestConnection.Contract)
	if mcnErr != nil {
		return InvalidConnectionContractError(pc, manifestConnection.Contract, moduleName)
	}
	ms, err := artifacts.ParseSource(found.ModuleSource)
	if err != nil {
		return InvalidModuleFormatError(pc, found.ModuleSource)
	}
	m, mErr := resolver.ResolveModule(ctx, *ms)
	if mErr != nil {
		return ModuleLookupFailedError(pc, found.ModuleSource, mErr)
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
			return MismatchedConnectionContractError(pc, found.Name, manifestConnection.Contract)
		}
	}

	return nil
}
