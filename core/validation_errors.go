package core

import (
	"fmt"
	"github.com/BSick7/go-api/errors"
	"github.com/nullstone-io/module/config"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

func InvalidYamlError(filename string, err error) errors.ValidationError {
	return errors.ValidationError{
		Context: fmt.Sprintf("file: %s", filename),
		Message: fmt.Sprintf("invalid .nullstone/%s file: %s", filename, err),
	}
}

func VariableDoesNotExistError(path, name, moduleName string) errors.ValidationError {
	return errors.ValidationError{
		Context: fmt.Sprintf("%s.vars.%s", path, name),
		Message: fmt.Sprintf("variable does not exist on the module (%s)", moduleName),
	}
}

func ConnectionDoesNotExistError(path, name, moduleName string) errors.ValidationError {
	return errors.ValidationError{
		Context: fmt.Sprintf("%s.connections.%s", path, name),
		Message: fmt.Sprintf("connection does not exist on the module (%s)", moduleName),
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

func MismatchedConnectionContractError(path string, blockConfig BlockConfiguration, connection config.Connection) errors.ValidationError {
	return errors.ValidationError{
		Context: path,
		Message: fmt.Sprintf("block (%s) does not match the required contract (%s) for the capability connection", blockConfig.Name, connection.Contract),
	}
}

func UnsupportedAppCategoryError(path, moduleSource, subcategory string) errors.ValidationError {
	return errors.ValidationError{
		Context: fmt.Sprintf("%s.module", path),
		Message: fmt.Sprintf("module (%s) does not support application category (%s)", moduleSource, subcategory),
	}
}
