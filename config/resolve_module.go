package config

import (
	"context"
	"github.com/BSick7/go-api/errors"
	"github.com/nullstone-io/iac/core"
	"gopkg.in/nullstone-io/go-api-client.v0/artifacts"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
	"strings"
)

func ResolveModule(ctx context.Context, resolver core.ValidateResolver, ic core.IacContext, pc core.ObjectPathContext,
	source, version string, contract types.ModuleContractName) (*types.Module, *types.ModuleVersion, *errors.ValidationError) {

	if source == "" {
		return nil, nil, RequiredModuleError(ic, pc)
	}

	ms, err := artifacts.ParseSource(source)
	if err != nil {
		return nil, nil, InvalidModuleFormatError(ic, pc, source)
	}
	m, mv, err := resolver.ResolveModuleVersion(ctx, *ms, version)
	if err != nil {
		return nil, nil, ModuleVersionLookupFailedError(ic, pc, source, version, err)
	}
	if m == nil {
		return nil, nil, MissingModuleError(ic, pc, source)
	}

	mcn2 := types.ModuleContractName{
		Category:    string(m.Category),
		Subcategory: string(m.Subcategory),
		Provider:    strings.Join(m.ProviderTypes, ","),
		Platform:    m.Platform,
		Subplatform: m.Subplatform,
	}
	if ok := contract.Match(mcn2); !ok {
		return nil, nil, InvalidModuleContractError(ic, pc, source, contract, mcn2)
	}

	if mv == nil {
		return nil, nil, MissingModuleVersionError(ic, pc, ms.String(), version)
	}

	return m, mv, nil
}
