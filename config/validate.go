package config

import (
	"context"
	errs "errors"
	"fmt"
	"github.com/BSick7/go-api/errors"
	"github.com/nullstone-io/iac/core"
	"github.com/nullstone-io/module/config"
	"gopkg.in/nullstone-io/go-api-client.v0/artifacts"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
	"strings"
)

func ValidateVariables(repoName, filename, path string, variables map[string]any, expectedVariables map[string]config.Variable, moduleName string) errors.ValidationErrors {
	ve := errors.ValidationErrors{}
	for k, _ := range variables {
		if _, ok := expectedVariables[k]; !ok {
			ve = append(ve, core.VariableDoesNotExistError(repoName, filename, path, k, moduleName))
		}
	}
	if len(ve) > 0 {
		return ve
	}
	return nil
}

// ValidateConnections performs validation on all IaC connections by matching them against connections in the module
func ValidateConnections(ctx context.Context, resolver *find.ResourceResolver, repoName, filename, path string, connections types.ConnectionTargets, manifestConnections map[string]config.Connection, moduleName string) error {
	ve := errors.ValidationErrors{}
	for key, conn := range connections {
		conPath := fmt.Sprintf("%s.connections.%s", path, key)
		manifestConnection, found := manifestConnections[key]
		if !found {
			ve = append(ve, core.ConnectionDoesNotExistError(repoName, filename, path, key, moduleName))
			continue
		}
		err := ValidateConnection(ctx, resolver, repoName, filename, conPath, key, conn, manifestConnection, moduleName)
		if err != nil {
			var verrs errors.ValidationErrors
			if errs.As(err, &verrs) {
				ve = append(ve, verrs...)
			}
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
func ValidateConnection(ctx context.Context, resolver *find.ResourceResolver, repoName, filename, path string, connName string, connection types.ConnectionTarget, manifestConnection config.Connection, moduleName string) error {
	if connection.BlockName == "" {
		return errors.ValidationErrors{core.MissingConnectionBlockError(repoName, filename, path)}
	}

	found, err := resolver.FindBlock(ctx, connection)
	if err != nil {
		if find.IsMissingResource(err) {
			return errors.ValidationErrors{core.MissingConnectionTargetError(repoName, filename, path, err)}
		}
		return err
	}

	mcn1, mcnErr := types.ParseModuleContractName(manifestConnection.Contract)
	if mcnErr != nil {
		return errors.ValidationErrors{core.InvalidConnectionContractError(repoName, filename, path, connName, manifestConnection.Contract, moduleName)}
	}
	ms, err := artifacts.ParseSource(found.ModuleSource)
	if err != nil {
		return errors.ValidationErrors{core.InvalidModuleFormatError(repoName, filename, path, found.ModuleSource)}
	}
	m, mErr := resolver.ApiClient.Modules().Get(ctx, ms.OrgName, ms.ModuleName)
	if mErr != nil {
		return fmt.Errorf("module lookup failed (%s): %w", found.ModuleSource, mErr)
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
			return errors.ValidationErrors{core.MismatchedConnectionContractError(repoName, filename, path, found.Name, manifestConnection.Contract)}
		}
	}

	return nil
}

func hasInvalidChars(r rune) bool {
	return (r < 'A' || r > 'z') && r != '_' && (r < '0' || r > '9')
}

func startsWithNumber(s string) bool {
	return s[0] >= '0' && s[0] <= '9'
}

func ValidateEnvVariables(repoName, filename, path string, envVariables map[string]string) error {
	ve := errors.ValidationErrors{}

	for k, _ := range envVariables {
		if startsWithNumber(k) {
			ve = append(ve, core.EnvVariableKeyStartsWithNumberError(repoName, filename, path, k))
		}
		if strings.IndexFunc(k, hasInvalidChars) != -1 {
			ve = append(ve, core.EnvVariableKeyInvalidCharsError(repoName, filename, path, k))
		}
	}

	if len(ve) > 0 {
		return ve
	}
	return nil
}

// ValidateCapabilities performs validation on a all IaC capabilities within an application
func ValidateCapabilities(ctx context.Context, resolver *find.ResourceResolver, repoName, filename, path string, capabilities CapabilityConfigurations, subcategory types.SubcategoryName) error {
	ve := errors.ValidationErrors{}
	for i, iacCap := range capabilities {
		capPath := fmt.Sprintf("%s.capabilities[%d]", path, i)
		err := ValidateCapability(ctx, resolver, repoName, filename, capPath, iacCap, string(subcategory))
		if err != nil {
			var verrs errors.ValidationErrors
			if errs.As(err, &verrs) {
				ve = append(ve, verrs...)
			}
		}
	}

	if len(ve) > 0 {
		return ve
	}
	return nil
}

func ValidateCapability(ctx context.Context, resolver *find.ResourceResolver, repoName, filename, path string, iacCap CapabilityConfiguration, subcategory string) error {
	// ensure the module is a capability module and supports the provider type (e.g. aws, gcp, azure)
	providerType, err := resolver.ResolveCurProviderType(ctx)
	if err != nil {
		return fmt.Errorf("unable to resolve current provider type: %w", err)
	}
	contract := fmt.Sprintf("capability/%s/*", providerType)
	m, mv, err := ResolveModule(ctx, resolver, repoName, filename, path, iacCap.ModuleSource, iacCap.ModuleSourceVersion, contract)
	if err != nil {
		return err
	}

	ve := errors.ValidationErrors{}
	// check to make sure the capability module supports the subcategory
	// examples are "container", "serverless", "static-site", "server"
	if m != nil {
		found := false
		for _, cat := range m.AppCategories {
			if cat == subcategory {
				found = true
				break
			}
		}
		if !found {
			ve = append(ve, core.UnsupportedAppCategoryError(repoName, filename, path, iacCap.ModuleSource, subcategory))
		}
	}

	// if we were able to find the module version
	//   1. validate each of the variables to ensure the module supports them
	//   2. validate each of the connections to ensure the block matches the connection contract
	if mv != nil {
		moduleName := fmt.Sprintf("%s@%s", iacCap.ModuleSource, iacCap.ModuleSourceVersion)
		verrs := ValidateVariables(repoName, filename, path, iacCap.Variables, mv.Manifest.Variables, moduleName)
		ve = append(ve, verrs...)

		err := ValidateConnections(ctx, resolver, repoName, filename, path, iacCap.Connections, mv.Manifest.Connections, moduleName)
		if err != nil {
			var verrs errors.ValidationErrors
			if errs.As(err, &verrs) {
				ve = append(ve, verrs...)
			}
		}
	}

	if len(ve) > 0 {
		return ve
	}
	return nil
}

func ValidateBlock(ctx context.Context, resolver *find.ResourceResolver, repoName, filename, yamlPath, contract, moduleSource, moduleSourceVersion string, variables map[string]any, connections types.ConnectionTargets, envVars map[string]string, capabilities CapabilityConfigurations) error {
	m, mv, err := ResolveModule(ctx, resolver, repoName, filename, yamlPath, moduleSource, moduleSourceVersion, contract)
	if err != nil {
		return err
	}

	moduleName := fmt.Sprintf("%s/%s@%s", m.OrgName, m.Name, mv.Version)

	ve := errors.ValidationErrors{}
	ve = append(ve, ValidateVariables(repoName, filename, yamlPath, variables, mv.Manifest.Variables, moduleName)...)

	if connections != nil {
		err := ValidateConnections(ctx, resolver, repoName, filename, yamlPath, connections, mv.Manifest.Connections, moduleName)
		if err != nil {
			var verrs errors.ValidationErrors
			if errs.As(err, &verrs) {
				ve = append(ve, verrs...)
			}
		}
	}

	if envVars != nil {
		err := ValidateEnvVariables(repoName, filename, yamlPath, envVars)
		if err != nil {
			var verrs errors.ValidationErrors
			if errs.As(err, &verrs) {
				ve = append(ve, verrs...)
			}
		}
	}

	if capabilities != nil {
		err := ValidateCapabilities(ctx, resolver, repoName, filename, yamlPath, capabilities, m.Subcategory)
		if err != nil {
			var verrs errors.ValidationErrors
			if errs.As(err, &verrs) {
				ve = append(ve, verrs...)
			}
		}

	}

	if len(ve) > 0 {
		return ve
	}
	return nil
}
