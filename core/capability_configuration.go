package core

import (
	"fmt"
	"github.com/BSick7/go-api/errors"
	"github.com/nullstone-io/module/config"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
	"log"
	"strings"
)

type CapabilityConfiguration struct {
	ModuleSource        string            `yaml:"module" json:"module"`
	ModuleSourceVersion *string           `yaml:"module_version" json:"moduleVersion"`
	Variables           map[string]any    `yaml:"vars" json:"vars"`
	Connections         ConnectionTargets `yaml:"connections" json:"connections"`
	Namespace           *string           `yaml:"namespace" json:"namespace"`
}

type InvalidConfigurationError struct {
	Err error
}

func (e InvalidConfigurationError) Error() string {
	return fmt.Sprintf("invalid app configuration: %s", e.Err.Error())
}

func (c CapabilityConfiguration) Normalize(resolver *find.ResourceResolver) (CapabilityConfiguration, error) {
	if err := NormalizeConnectionTargets(c.Connections, resolver); err != nil {
		return c, err
	}
	return c, nil
}

// validateConnections loops through each of the connections in the configuration
//  1. ensures that the block exists
//  2. ensures that the module for the block matches the connection contract
func (c CapabilityConfiguration) validateConnections(path string, connections map[string]config.Connection, resolver *find.ResourceResolver) (errors.ValidationErrors, error) {
	ve := errors.ValidationErrors{}
	for key, conn := range c.Connections {
		conPath := fmt.Sprintf("%s.connections.%s", path, key)
		connection, found := connections[key]
		if !found {
			ve = append(ve, errors.ValidationError{Context: conPath, Message: fmt.Sprintf("connection does not exist on the module (%s@%s)", c.ModuleSource, *c.ModuleSourceVersion)})
			continue
		}
		block, err := resolver.FindBlock(types.ConnectionTarget(conn))
		if err != nil {
			if find.IsMissingResource(err) {
				ve = append(ve, errors.ValidationError{Context: conPath, Message: fmt.Sprintf("connection is invalid, %s", err)})
				continue
			}
			return ve, err
		}

		mcn1, mcnErr := types.ParseModuleContractName(connection.Contract)
		if mcnErr != nil {
			log.Printf("unable to validate capabilility (%s@%s) connection (%s): connection contract name (%s) parse failed: %s\n", c.ModuleSource, *c.ModuleSourceVersion, key, connection.Contract, mcnErr)
			ve = append(ve, errors.ValidationError{Context: conPath, Message: fmt.Sprintf("an error occurred verifying contract: connection contract (%s) has an incorrect format", connection.Contract)})
		}
		parts := strings.Split(block.ModuleSource, "/")
		if len(parts) != 2 {
			ve = append(ve, errors.ValidationError{Context: fmt.Sprintf("%s.module", path), Message: fmt.Sprintf("module (%s) must be in the format \"org/name\"", block.ModuleSource)})
			continue
		}
		m, mErr := resolver.ApiClient.Modules().Get(parts[0], parts[1])
		if mErr != nil {
			return ve, fmt.Errorf("unable to validate capability (%s@%s) connection (%s): module lookup failed (%s): %w", c.ModuleSource, *c.ModuleSourceVersion, key, block.ModuleSource, mErr)
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
				ve = append(ve, errors.ValidationError{Context: conPath, Message: fmt.Sprintf("block (%s) does not match the required contract (%s) for the capability connection", block.Name, connection.Contract)})
			}
		}
	}

	return ve, nil
}
