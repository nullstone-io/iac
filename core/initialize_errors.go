package core

import (
	"fmt"

	"github.com/BSick7/go-api/errors"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

var (
	_ error = InitializeError{}
)

type InitializeError struct {
	IacContext        IacContext        `json:"iacContext"`
	ObjectPathContext ObjectPathContext `json:"objectPathContext"`
	ErrorMessage      string            `json:"errorMessage"`
}

func (e InitializeError) Error() string {
	return fmt.Sprintf("%s => %s", e.IacContext.Context(e.ObjectPathContext), e.ErrorMessage)
}

func (e InitializeError) ToValidationError() errors.ValidationError {
	return errors.ValidationError{
		Context: e.IacContext.Context(e.ObjectPathContext),
		Message: e.ErrorMessage,
	}
}

type InitializeErrors []InitializeError

func (s InitializeErrors) ToValidationErrors() errors.ValidationErrors {
	if len(s) == 0 {
		return nil
	}
	ve := errors.ValidationErrors{}
	for _, re := range s {
		ve = append(ve, re.ToValidationError())
	}
	return ve
}

func RequiredModuleError(pc ObjectPathContext) *InitializeError {
	return &InitializeError{
		ObjectPathContext: pc,
		ErrorMessage:      "Module is required",
	}
}

func InvalidResolveModuleFormatError(pc ObjectPathContext, moduleSource string) *InitializeError {
	return &InitializeError{
		ObjectPathContext: pc.SubField("module"),
		ErrorMessage:      fmt.Sprintf("Invalid module format (%s) - must be in the format \"<module-org>/<module-name>\"", moduleSource),
	}
}

func ModuleVersionLookupFailedError(pc ObjectPathContext, source, version string, err error) *InitializeError {
	return &InitializeError{
		ObjectPathContext: pc.SubField("module_version"),
		ErrorMessage:      fmt.Sprintf("Module version (%s@%s) lookup failed: %s", source, version, err),
	}
}

func MissingModuleError(pc ObjectPathContext, moduleSource string) *InitializeError {
	return &InitializeError{
		ObjectPathContext: pc.SubField("module"),
		ErrorMessage:      fmt.Sprintf("Module (%s) does not exist", moduleSource),
	}
}

func InvalidModuleContractError(pc ObjectPathContext, moduleSource string, want, got types.ModuleContractName) *InitializeError {
	return &InitializeError{
		ObjectPathContext: pc.SubField("module"),
		ErrorMessage:      fmt.Sprintf("Module (%s) must be %s module and match the contract (%s), it is defined as %s", moduleSource, want.Category, want, got),
	}
}

func MissingModuleVersionError(pc ObjectPathContext, source, version string) *InitializeError {
	return &InitializeError{
		ObjectPathContext: pc.SubField("module_version"),
		ErrorMessage:      fmt.Sprintf("Module version (%s@%s) does not exist", source, version),
	}
}

func MissingRequiredConnectionError(pc ObjectPathContext, connName string) InitializeError {
	return InitializeError{
		ObjectPathContext: pc,
		ErrorMessage:      fmt.Sprintf("Connection (%s) is required", connName),
	}
}
