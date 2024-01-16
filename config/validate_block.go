package config

import (
	errs "errors"
	"fmt"
	"github.com/BSick7/go-api/errors"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

func ValidateBlock(resolver *find.ResourceResolver, configBlocks []BlockConfiguration, yamlPath, contract, moduleSource, moduleSourceVersion string, variables map[string]any, connections types.ConnectionTargets, capabilities CapabilityConfigurations) error {
	m, mv, err := ResolveModule(resolver, yamlPath, moduleSource, moduleSourceVersion, contract)
	if err != nil {
		return err
	}

	moduleName := fmt.Sprintf("%s/%s@%s", m.OrgName, m.Name, mv.Version)

	ve := errors.ValidationErrors{}
	ve = append(ve, ValidateVariables(yamlPath, variables, mv.Manifest.Variables, moduleName)...)

	if connections != nil {
		err := ValidateConnections(resolver, configBlocks, yamlPath, connections, mv.Manifest.Connections, moduleName)
		if err != nil {
			var verrs errors.ValidationErrors
			if errs.As(err, &verrs) {
				ve = append(ve, verrs...)
			}
		}
	}

	if capabilities != nil {
		err := ValidateCapabilities(resolver, configBlocks, yamlPath, capabilities, m.Subcategory)
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
