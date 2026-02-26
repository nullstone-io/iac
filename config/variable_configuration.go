package config

import (
	"strings"

	"github.com/nullstone-io/iac/core"
	"github.com/nullstone-io/module/config"
)

type VariableConfigurations map[string]*VariableConfiguration

func (s VariableConfigurations) Initialize(blockManifest config.Manifest) core.InitializeErrors {
	if len(s) == 0 {
		return nil
	}
	errs := core.InitializeErrors{}
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
	if !isValueCompatibleWithType(c.Value, c.Schema.Type) {
		return core.VariableIncompatibleTypeError(pc, c.Schema.Type, c.Value)
	}
	return nil
}

// isValueCompatibleWithType checks whether the given value (parsed from YAML) is compatible
// with the declared Terraform variable type. The type string may be a simple type (e.g. "string",
// "number", "bool") or a complex type expression (e.g. "list(string)", "map(number)",
// "object({...})", "set(string)", "tuple([...])"). Only the base type name is used for the check.
func isValueCompatibleWithType(value any, varType string) bool {
	if value == nil {
		return true
	}

	baseType := extractBaseType(varType)

	switch baseType {
	case "string":
		_, ok := value.(string)
		return ok
	case "number":
		return isNumeric(value)
	case "bool":
		_, ok := value.(bool)
		return ok
	case "list", "set", "tuple":
		_, ok := value.([]any)
		return ok
	case "map", "object":
		_, ok := value.(map[string]any)
		return ok
	default:
		// Unknown type — allow any value to avoid false positives
		return true
	}
}

// extractBaseType returns the base type name from a Terraform type expression.
// Examples:
//
//	"string"                          → "string"
//	"list(string)"                    → "list"
//	"map(number)"                     → "map"
//	"object({ enabled : bool ... })"  → "object"
//	"set(string)"                     → "set"
//	"tuple([string, number])"         → "tuple"
func extractBaseType(varType string) string {
	varType = strings.TrimSpace(varType)
	if idx := strings.IndexAny(varType, "(["); idx != -1 {
		return strings.TrimSpace(varType[:idx])
	}
	return varType
}

// isNumeric returns true for any Go numeric type that YAML or JSON might produce.
func isNumeric(value any) bool {
	switch value.(type) {
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		return true
	}
	return false
}
