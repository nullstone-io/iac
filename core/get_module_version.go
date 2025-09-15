package core

import (
	"context"
	"strings"

	"gopkg.in/nullstone-io/go-api-client.v0/artifacts"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

func GetModuleVersion(ctx context.Context, resolver ModuleVersionResolver, pc ObjectPathContext,
	source, version string, contract types.ModuleContractName) (*types.Module, *types.ModuleVersion, *InitializeError) {

	if source == "" {
		return nil, nil, RequiredModuleError(pc)
	}

	ms, err := artifacts.ParseSource(source)
	if err != nil {
		return nil, nil, InvalidResolveModuleFormatError(pc, source)
	}
	m, mv, err := resolver.ResolveModuleVersion(ctx, *ms, version)
	if err != nil {
		return nil, nil, ModuleVersionLookupFailedError(pc, source, version, err)
	}
	if m == nil {
		return nil, nil, MissingModuleError(pc, source)
	}

	mcn2 := types.ModuleContractName{
		Category:    string(m.Category),
		Subcategory: string(m.Subcategory),
		Provider:    strings.Join(m.ProviderTypes, ","),
		Platform:    m.Platform,
		Subplatform: m.Subplatform,
	}
	if ok := contract.Match(mcn2); !ok {
		return nil, nil, InvalidModuleContractError(pc, source, contract, mcn2)
	}

	if mv == nil {
		return nil, nil, MissingModuleVersionError(pc, ms.String(), version)
	}

	return m, mv, nil
}
