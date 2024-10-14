package core

import (
	"fmt"
	"github.com/BSick7/go-api/errors"
)

var (
	_ error = ValidateError{}
)

type ValidateError struct {
	IacContext        IacContext        `json:"iacContext"`
	ObjectPathContext ObjectPathContext `json:"objectPathContext"`
	ErrorMessage      string            `json:"errorMessage"`
}

func (e ValidateError) Error() string {
	return fmt.Sprintf("%s => %s", e.IacContext.Context(e.ObjectPathContext), e.ErrorMessage)
}

func (e ValidateError) ToValidationError() errors.ValidationError {
	return errors.ValidationError{
		Context: e.IacContext.Context(e.ObjectPathContext),
		Message: e.ErrorMessage,
	}
}

type ValidateErrors []ValidateError

func (s ValidateErrors) ToValidationErrors() errors.ValidationErrors {
	if len(s) == 0 {
		return nil
	}
	ve := errors.ValidationErrors{}
	for _, re := range s {
		ve = append(ve, re.ToValidationError())
	}
	return ve
}

func VariableDoesNotExistError(pc ObjectPathContext, moduleName string) *ValidateError {
	return &ValidateError{
		ObjectPathContext: pc,
		ErrorMessage:      fmt.Sprintf("Variable does not exist on the module (%s)", moduleName),
	}
}

func ConnectionDoesNotExistError(pc ObjectPathContext, moduleName string) *ValidateError {
	return &ValidateError{
		ObjectPathContext: pc,
		ErrorMessage:      fmt.Sprintf("Connection does not exist on the module (%s)", moduleName),
	}
}

func MissingConnectionBlockError(pc ObjectPathContext) *ValidateError {
	return &ValidateError{
		ObjectPathContext: pc,
		ErrorMessage:      fmt.Sprintf("Connection must have a block_name to identify which block it is connected to"),
	}
}

func InvalidConnectionContractError(pc ObjectPathContext, contract, moduleName string) *ValidateError {
	return &ValidateError{
		ObjectPathContext: pc,
		ErrorMessage:      fmt.Sprintf("Connection contract (contract=%s) in module (%s) is invalid", contract, moduleName),
	}
}

func MismatchedConnectionContractError(pc ObjectPathContext, blockName, connectionContract string) *ValidateError {
	return &ValidateError{
		ObjectPathContext: pc,
		ErrorMessage:      fmt.Sprintf("Block (%s) does not match the required contract (%s) for the capability connection", blockName, connectionContract),
	}
}
