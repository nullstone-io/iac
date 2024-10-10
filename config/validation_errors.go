package config

import (
	errors2 "errors"
	"fmt"
	"github.com/BSick7/go-api/errors"
	"github.com/nullstone-io/iac/core"
)

func NewValidationError(ic core.IacContext, ipc core.ObjectPathContext, msg string) *errors.ValidationError {
	return &errors.ValidationError{
		Context: ic.Context(ipc),
		Message: msg,
	}
}

func VariableDoesNotExistError(ic core.IacContext, ipc core.ObjectPathContext, moduleName string) errors.ValidationError {
	err := NewValidationError(ic, ipc, fmt.Sprintf("Variable does not exist on the module (%s)", moduleName))
	return *err
}

func EnvVariableKeyStartsWithNumberError(ic core.IacContext, pc core.ObjectPathContext) errors.ValidationError {
	err := NewValidationError(ic, pc, "Invalid environment variable, key must not start with a number")
	return *err
}

func EnvVariableKeyInvalidCharsError(ic core.IacContext, pc core.ObjectPathContext) errors.ValidationError {
	err := NewValidationError(ic, pc, "Invalid environment variable, key must contain only letters, numbers, and underscores")
	return *err
}

func ConnectionDoesNotExistError(ic core.IacContext, pc core.ObjectPathContext, moduleName string) errors.ValidationError {
	err := NewValidationError(ic, pc, fmt.Sprintf("Connection does not exist on the module (%s)", moduleName))
	return *err
}

func MissingConnectionBlockError(ic core.IacContext, pc core.ObjectPathContext) *errors.ValidationError {
	return NewValidationError(ic, pc, fmt.Sprintf("Connection must have a block_name to identify which block it is connected to"))
}

func MissingConnectionTargetError(ic core.IacContext, pc core.ObjectPathContext, err error) *errors.ValidationError {
	return NewValidationError(ic, pc, fmt.Sprintf("Connection is invalid, %s", err))
}

func LookupConnectionTargetFailedError(ic core.IacContext, pc core.ObjectPathContext, err error) *errors.ValidationError {
	return NewValidationError(ic, pc, fmt.Sprintf("Failed to validate connection, error when looking up connection target: %s", err))
}

func ModuleLookupFailedError(ic core.IacContext, pc core.ObjectPathContext, moduleSource string, err error) *errors.ValidationError {
	return NewValidationError(ic, pc.SubField("module"), fmt.Sprintf("Module (%s) lookup failed: %s", moduleSource, err))
}

func InvalidModuleFormatError(ic core.IacContext, pc core.ObjectPathContext, moduleSource string) *errors.ValidationError {
	return NewValidationError(ic, pc.SubField("module"), fmt.Sprintf("Invalid module format (%s) - must be in the format \"<module-org>/<module-name>\"", moduleSource))
}

func InvalidModuleContractParserError(ic core.IacContext, pc core.ObjectPathContext, moduleSource, contract string, err error) *errors.ValidationError {
	return NewValidationError(ic, pc.SubField("module"), fmt.Sprintf("Invalid module (%s) contract (%s), parse failed: %s", moduleSource, contract, err))
}

func InvalidConnectionContractError(ic core.IacContext, pc core.ObjectPathContext, contract, moduleName string) *errors.ValidationError {
	return NewValidationError(ic, pc, fmt.Sprintf("Connection contract (contract=%s) in module (%s) is invalid", contract, moduleName))
}

func MismatchedConnectionContractError(ic core.IacContext, pc core.ObjectPathContext, blockName, connectionContract string) *errors.ValidationError {
	return NewValidationError(ic, pc, fmt.Sprintf("Block (%s) does not match the required contract (%s) for the capability connection", blockName, connectionContract))
}

func LookupProviderTypeFailedError(ic core.IacContext, pc core.ObjectPathContext, err error) errors.ValidationError {
	verr := NewValidationError(ic, pc, fmt.Sprintf("Lookup for capability provider type failed: %s", err))
	return *verr
}

func UnsupportedAppCategoryError(ic core.IacContext, pc core.ObjectPathContext, moduleSource, subcategory string) errors.ValidationError {
	err := NewValidationError(ic, pc, fmt.Sprintf("Module (%s) does not support application category (%s)", moduleSource, subcategory))
	return *err
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
