package config

import (
	"context"
	"github.com/BSick7/go-api/errors"
	"github.com/nullstone-io/iac/core"
	"gopkg.in/nullstone-io/go-api-client.v0/artifacts"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
	"strings"
)

func ResolveModule(ctx context.Context, resolver *find.ResourceResolver, ic core.IacContext, pc core.YamlPathContext, moduleSource, moduleSourceVersion, contract string) (*types.Module, *types.ModuleVersion, *errors.ValidationError) {
	if moduleSource == "" {
		return nil, nil, RequiredModuleError(ic, pc)
	}

	ms, err := artifacts.ParseSource(moduleSource)
	if err != nil {
		return nil, nil, InvalidModuleFormatError(ic, pc, moduleSource)
	}
	// TODO: Add support for ms.Host
	m, err := resolver.ApiClient.Modules().Get(ctx, ms.OrgName, ms.ModuleName)
	if err != nil {
		return nil, nil, ModuleLookupFailedError(ic, pc, moduleSource, err)
	}
	if m == nil {
		return nil, nil, MissingModuleError(ic, pc, moduleSource)
	}
	mcn1, err := types.ParseModuleContractName(contract)
	if err != nil {
		return nil, nil, InvalidModuleContractParserError(ic, pc, moduleSource, contract, err)
	}
	mcn2 := types.ModuleContractName{
		Category:    string(m.Category),
		Subcategory: string(m.Subcategory),
		Provider:    strings.Join(m.ProviderTypes, ","),
		Platform:    m.Platform,
		Subplatform: m.Subplatform,
	}
	if ok := mcn1.Match(mcn2); !ok {
		return nil, nil, InvalidModuleContractError(ic, pc, moduleSource, mcn1, mcn2)
	}

	var mv *types.ModuleVersion
	if moduleSourceVersion == "latest" {
		mv = m.LatestVersion
	} else {
		mv, err = resolver.ApiClient.ModuleVersions().Get(ctx, ms.OrgName, ms.ModuleName, moduleSourceVersion)
		if err != nil {
			return nil, nil, ModuleVersionLookupFailedError(ic, pc, moduleSource, moduleSourceVersion, err)
		}
	}

	if mv == nil {
		return nil, nil, MissingModuleVersionError(ic, pc, ms.String(), moduleSourceVersion)
	}

	return m, mv, nil
}
