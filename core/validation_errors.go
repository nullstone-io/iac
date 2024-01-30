package core

import (
	errors2 "errors"
	"fmt"
	"github.com/BSick7/go-api/errors"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

func VariableDoesNotExistError(repoName, filename, path, name, moduleName string) errors.ValidationError {
	return errors.ValidationError{
		Context: fmt.Sprintf("%s#%s (%s.vars.%s)\n", repoName, filename, path, name),
		Message: fmt.Sprintf("Variable does not exist on the module (%s)", moduleName),
	}
}

func EnvVariableKeyStartsWithNumberError(repoName, filename, path, key string) errors.ValidationError {
	return errors.ValidationError{
		Context: fmt.Sprintf("%s#%s (%s.env_vars.%s)\n", repoName, filename, path, key),
		Message: fmt.Sprintf("Invalid environment variable key (%s) - it must not start with a number", key),
	}
}

func EnvVariableKeyInvalidCharsError(repoName, filename, path, key string) errors.ValidationError {
	return errors.ValidationError{
		Context: fmt.Sprintf("%s#%s (%s.env_vars.%s)\n", repoName, filename, path, key),
		Message: fmt.Sprintf("Invalid environment variable key (%s) - it must only contain letters, numbers, and underscores", key),
	}
}

func ConnectionDoesNotExistError(repoName, filename, path, name, moduleName string) errors.ValidationError {
	return errors.ValidationError{
		Context: fmt.Sprintf("%s#%s (%s.connections.%s)\n", repoName, filename, path, name),
		Message: fmt.Sprintf("Connection does not exist on the module (%s)", moduleName),
	}
}

func MissingConnectionBlockError(repoName, filename, path string) errors.ValidationError {
	return errors.ValidationError{
		Context: fmt.Sprintf("%s#%s (%s.block_name)\n", repoName, filename, path),
		Message: fmt.Sprintf("Connection must have a block_name to identify which block it is connected to"),
	}
}

func MissingConnectionTargetError(repoName, filename, path string, err error) errors.ValidationError {
	return errors.ValidationError{
		Context: fmt.Sprintf("%s#%s (%s)\n", repoName, filename, path),
		Message: fmt.Sprintf("Connection is invalid, %s", err),
	}
}

func InvalidModuleFormatError(repoName, filename, path, moduleSource string) errors.ValidationError {
	return errors.ValidationError{
		Context: fmt.Sprintf("%s#%s (%s.module)\n", repoName, filename, path),
		Message: fmt.Sprintf("Invalid module format (%s) - must be in the format \"<module-org>/<module-name>\"", moduleSource),
	}
}

func RequiredModuleError(repoName, filename, path string) errors.ValidationError {
	return errors.ValidationError{
		Context: fmt.Sprintf("%s#%s (%s.module)\n", repoName, filename, path),
		Message: fmt.Sprintf("Module is required"),
	}
}

func MissingModuleError(repoName, filename, path, moduleSource string) errors.ValidationError {
	return errors.ValidationError{
		Context: fmt.Sprintf("%s#%s (%s.module)\n", repoName, filename, path),
		Message: fmt.Sprintf("Module (%s) does not exist", moduleSource),
	}
}

func InvalidModuleContractError(repoName, filename, path, moduleSource string, want, got types.ModuleContractName) errors.ValidationError {
	return errors.ValidationError{
		Context: fmt.Sprintf("%s#%s (%s.module)\n", repoName, filename, path),
		Message: fmt.Sprintf("Module (%s) must be %s module and match the contract (%s), it is defined as %s", moduleSource, want.Category, want, got),
	}
}

func MissingModuleVersionError(repoName, filename, path, source, version string) errors.ValidationError {
	return errors.ValidationError{
		Context: fmt.Sprintf("%s#%s (%s.module_version)\n", repoName, filename, path),
		Message: fmt.Sprintf("Module version (%s@%s) does not exist", source, version),
	}
}

func InvalidConnectionContractError(repoName, filename, path, connName, contract, moduleName string) errors.ValidationError {
	return errors.ValidationError{
		Context: fmt.Sprintf("%s#%s (%s)\n", repoName, filename, path),
		Message: fmt.Sprintf("Connection contract (name=%s, contract=%s) in module (%s) is invalid", connName, contract, moduleName),
	}
}

func MismatchedConnectionContractError(repoName, filename, path string, blockName, connectionContract string) errors.ValidationError {
	return errors.ValidationError{
		Context: fmt.Sprintf("%s#%s (%s)\n", repoName, filename, path),
		Message: fmt.Sprintf("Block (%s) does not match the required contract (%s) for the capability connection", blockName, connectionContract),
	}
}

func UnsupportedAppCategoryError(repoName, filename, path, moduleSource, subcategory string) errors.ValidationError {
	return errors.ValidationError{
		Context: fmt.Sprintf("%s#%s (%s.module)\n", repoName, filename, path),
		Message: fmt.Sprintf("Module (%s) does not support application category (%s)", moduleSource, subcategory),
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
