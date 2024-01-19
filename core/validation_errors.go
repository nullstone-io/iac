package core

import (
	errors2 "errors"
	"fmt"
	"github.com/BSick7/go-api/errors"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

func VariableDoesNotExistError(path, name, moduleName string) errors.ValidationError {
	return errors.ValidationError{
		Context: fmt.Sprintf("%s.vars.%s", path, name),
		Message: fmt.Sprintf("variable does not exist on the module (%s)", moduleName),
	}
}

func EnvVariableKeyStartsWithNumberError(path, key string) errors.ValidationError {
	return errors.ValidationError{
		Context: fmt.Sprintf("%s.env_vars.%s", path, key),
		Message: fmt.Sprintf("invalid environment variable key (%s) - it must not start with a number", key),
	}
}

func EnvVariableKeyInvalidCharsError(path, key string) errors.ValidationError {
	return errors.ValidationError{
		Context: fmt.Sprintf("%s.env_vars.%s", path, key),
		Message: fmt.Sprintf("invalid environment variable key (%s) - it must only contain letters, numbers, and underscores", key),
	}
}

func ConnectionDoesNotExistError(path, name, moduleName string) errors.ValidationError {
	return errors.ValidationError{
		Context: fmt.Sprintf("%s.connections.%s", path, name),
		Message: fmt.Sprintf("connection does not exist on the module (%s)", moduleName),
	}
}

func MissingConnectionBlockError(path string) errors.ValidationError {
	return errors.ValidationError{
		Context: fmt.Sprintf("%s.block_name", path),
		Message: fmt.Sprintf("connection must have a block_name to identify which block it is connected to"),
	}
}

func MissingConnectionTargetError(path string, err error) errors.ValidationError {
	return errors.ValidationError{
		Context: path,
		Message: fmt.Sprintf("connection is invalid, %s", err),
	}
}

func InvalidModuleFormatError(path, moduleSource string, err error) errors.ValidationError {
	return errors.ValidationError{
		Context: fmt.Sprintf("%s.module", path),
		Message: fmt.Sprintf("%s (%s), must be in the format \"<module-org>/<module-name>\"", err, moduleSource),
	}
}

func MissingModuleError(path, moduleSource string) errors.ValidationError {
	return errors.ValidationError{
		Context: fmt.Sprintf("%s.module", path),
		Message: fmt.Sprintf("module (%s) does not exist", moduleSource),
	}
}

func InvalidModuleContractError(path, moduleSource string, want, got types.ModuleContractName) errors.ValidationError {
	return errors.ValidationError{
		Context: fmt.Sprintf("%s.module", path),
		Message: fmt.Sprintf("module (%s) must be %s module and match the contract (%s), it is defined as %s", moduleSource, want.Category, want, got),
	}
}

func MissingModuleVersionError(path, source, version string) errors.ValidationError {
	return errors.ValidationError{
		Context: fmt.Sprintf("%s.module_version", path),
		Message: fmt.Sprintf("module version (%s@%s) does not exist", source, version),
	}
}

func InvalidConnectionContractError(path, connName, contract, moduleName string) errors.ValidationError {
	return errors.ValidationError{
		Context: path,
		Message: fmt.Sprintf("connection contract (name=%s, contract=%s) in module (%s) is invalid", connName, contract, moduleName),
	}
}

func MismatchedConnectionContractError(path string, blockName, connectionContract string) errors.ValidationError {
	return errors.ValidationError{
		Context: path,
		Message: fmt.Sprintf("block (%s) does not match the required contract (%s) for the capability connection", blockName, connectionContract),
	}
}

func UnsupportedAppCategoryError(path, moduleSource, subcategory string) errors.ValidationError {
	return errors.ValidationError{
		Context: fmt.Sprintf("%s.module", path),
		Message: fmt.Sprintf("module (%s) does not support application category (%s)", moduleSource, subcategory),
	}
}

// AsValidationErrors is a helper function to format an error into validation errors for the user to see
func AsValidationErrors(err error) errors.ValidationErrors {
	var verrs errors.ValidationErrors
	if errors2.As(err, &verrs) {
		return verrs
	}
	var verr errors.ValidationError
	if errors2.As(err, &verr) {
		return errors.ValidationErrors{verr}
	}
	return nil
}
