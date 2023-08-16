package core

import (
	"fmt"
	"github.com/BSick7/go-api/errors"
	"gopkg.in/nullstone-io/go-api-client.v0/artifacts"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
	"strings"
)

func ResolveModule(resolver *find.ResourceResolver, iacPath, moduleSource, moduleSourceVersion, contract string) (*types.Module, *types.ModuleVersion, errors.ValidationErrors, error) {
	if moduleSource == "" {
		return nil, nil, errors.ValidationErrors{{Context: iacPath, Message: "module is required"}}, nil
	}

	ms, err := artifacts.ParseSource(moduleSource)
	if err != nil {
		return nil, nil, errors.ValidationErrors{InvalidModuleFormatError(fmt.Sprintf("%s.module", iacPath), moduleSource, err)}, nil
	}
	// TODO: Add support for ms.Host
	m, err := resolver.ApiClient.Modules().Get(ms.OrgName, ms.ModuleName)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("unable to validate module (%s): module lookup failed: %w", moduleSource, err)
	}
	if m == nil {
		return nil, nil, errors.ValidationErrors{MissingModuleError(iacPath, moduleSource)}, nil
	}
	mcn1, err := types.ParseModuleContractName(contract)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("unable to validate module (%s): module contract name (%s) parse failed: %w", moduleSource, contract, err)
	}
	mcn2 := types.ModuleContractName{
		Category:    string(m.Category),
		Subcategory: string(m.Subcategory),
		Provider:    strings.Join(m.ProviderTypes, ","),
		Platform:    m.Platform,
		Subplatform: m.Subplatform,
	}
	if ok := mcn1.Match(mcn2); !ok {
		return nil, nil, errors.ValidationErrors{InvalidModuleContractError(iacPath, moduleSource, mcn1, mcn2)}, nil
	}

	var mv *types.ModuleVersion
	if moduleSourceVersion == "latest" {
		mv = m.LatestVersion
	} else {
		mv, err = resolver.ApiClient.ModuleVersions().Get(ms.OrgName, ms.ModuleName, moduleSourceVersion)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("unable to validate module@version (%s@%s): module version lookup failed: %w", moduleSource, moduleSourceVersion, err)
		}
	}
	if mv == nil {
		return nil, nil, errors.ValidationErrors{MissingModuleVersionError(iacPath, ms.String(), moduleSourceVersion)}, nil
	}

	return m, mv, nil, nil
}
