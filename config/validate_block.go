package config

import (
	"fmt"
	"github.com/BSick7/go-api/errors"
	"github.com/nullstone-io/iac"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
)

func ValidateBlock(resolver *find.ResourceResolver, yamlPath, contract, moduleSource, moduleSourceVersion string, variables map[string]any, connections iac.ConnectionTargets, capabilities CapabilityConfigurations) (errors.ValidationErrors, error) {
	m, mv, verrs, err := iac.ResolveModule(resolver, yamlPath, moduleSource, moduleSourceVersion, contract)
	if err != nil {
		return nil, err
	} else if len(verrs) > 0 {
		return verrs, nil
	}

	moduleName := fmt.Sprintf("%s/%s@%s", m.OrgName, m.Name, mv.Version)

	ve := errors.ValidationErrors{}
	ve = append(ve, iac.ValidateVariables(yamlPath, variables, mv.Manifest.Variables, moduleName)...)

	if connections != nil {
		verrs, err := iac.ValidateConnections(resolver, yamlPath, connections, mv.Manifest.Connections, moduleName)
		if err != nil {
			return ve, err
		}
		ve = append(ve, verrs...)
	}

	if capabilities != nil {
		verrs, err := ValidateCapabilities(resolver, yamlPath, capabilities, m.Subcategory)
		if err != nil {
			return ve, err
		}
		ve = append(ve, verrs...)

	}
	return ve, nil
}
