package config

import (
	"github.com/nullstone-io/iac/core"
	"github.com/nullstone-io/module/config"
)

type VariableConfigurations map[string]*VariableConfiguration

func (s VariableConfigurations) Resolve(blockManifest config.Manifest) core.ResolveErrors {
	if len(s) == 0 {
		return nil
	}
	errs := core.ResolveErrors{}
	for key, c := range s {
		if schema, ok := blockManifest.Variables[key]; ok {
			c.Schema = &schema
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (s VariableConfigurations) Validate(pc core.ObjectPathContext, moduleName string) core.ValidateErrors {
	if len(s) == 0 {
		return nil
	}
	errs := core.ValidateErrors{}
	for k, cur := range s {
		if err := cur.Validate(pc.SubKey("vars", k), moduleName); err != nil {
			errs = append(errs, *err)
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

type VariableConfiguration struct {
	Value  any              `json:"value"`
	Schema *config.Variable `json:"schema"`
}

func (c *VariableConfiguration) Validate(pc core.ObjectPathContext, moduleName string) *core.ValidateError {
	if c.Schema == nil {
		return core.VariableDoesNotExistError(pc, moduleName)
	}
	// TODO: Perform type->value checks
	return nil
}
