package config

import (
	"fmt"
	"github.com/BSick7/go-api/errors"
	"github.com/nullstone-io/iac/core"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

type AppConfiguration struct {
	Name                string                   `yaml:"-" json:"name"`
	ModuleSource        string                   `yaml:"module" json:"module"`
	ModuleSourceVersion *string                  `yaml:"module_version,omitempty" json:"moduleVersion"`
	Variables           map[string]any           `yaml:"vars" json:"vars"`
	Capabilities        CapabilityConfigurations `yaml:"capabilities" json:"capabilities"`
	EnvVariables        map[string]string        `yaml:"environment" json:"envVars"`
}

func (a AppConfiguration) GetCapabilities(orgName string, stackId, blockId, envId int64) ([]types.Capability, error) {
	caps := make([]types.Capability, len(a.Capabilities))
	for i, cap := range a.Capabilities {
		updateCap := types.Capability{
			OrgName:             orgName,
			AppId:               blockId,
			EnvId:               envId,
			ModuleSource:        cap.ModuleSource,
			ModuleSourceVersion: "latest",
			Connections:         map[string]types.ConnectionTarget{},
		}
		if cap.ModuleSourceVersion != nil {
			updateCap.ModuleSourceVersion = *cap.ModuleSourceVersion
		}
		if cap.Namespace != nil {
			updateCap.Namespace = *cap.Namespace
		}
		for key, conn := range cap.Connections {
			target := types.ConnectionTarget{}
			// each connection must have a block_name to identify which block it is connected to
			if conn.BlockName == "" {
				return nil, core.InvalidConfigurationError{fmt.Errorf("The connection (%s) must have a block_name to identify which block it is connected to.", key)}
			}
			target.BlockName = conn.BlockName
			// each connection must also have a stack_id
			if conn.StackId != 0 {
				target.StackId = conn.StackId
			} else {
				target.StackId = stackId
			}
			target.EnvId = conn.EnvId
			updateCap.Connections[key] = target
		}
		caps[i] = updateCap
	}
	return caps, nil
}

func (a *AppConfiguration) Normalize(resolver *find.ResourceResolver) error {
	return a.Capabilities.Normalize(resolver)
}

func (a AppConfiguration) Validate(resolver *find.ResourceResolver, configBlocks []core.BlockConfiguration) (errors.ValidationErrors, error) {
	yamlPath := fmt.Sprintf("apps.%s", a.Name)
	return ValidateBlock(resolver, configBlocks, yamlPath, "app/*/*", a.ModuleSource, *a.ModuleSourceVersion, a.Variables, nil, a.Capabilities)
}
