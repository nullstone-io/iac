package config

import (
	"fmt"
	"github.com/nullstone-io/iac/core"
)

func EnvVariableKeyStartsWithNumberError(pc core.ObjectPathContext) core.ValidateError {
	return core.ValidateError{
		ObjectPathContext: pc,
		ErrorMessage:      "Invalid environment variable, key must not start with a number",
	}
}

func EnvVariableKeyInvalidCharsError(pc core.ObjectPathContext) core.ValidateError {
	return core.ValidateError{
		ObjectPathContext: pc,
		ErrorMessage:      "Invalid environment variable, key must contain only letters, numbers, and underscores",
	}
}

func UnsupportedAppCategoryError(pc core.ObjectPathContext, moduleSource, subcategory string) core.ValidateError {
	return core.ValidateError{
		ObjectPathContext: pc,
		ErrorMessage:      fmt.Sprintf("Module (%s) does not support application category (%s)", moduleSource, subcategory),
	}
}
