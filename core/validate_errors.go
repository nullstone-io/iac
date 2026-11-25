package core

import (
	"fmt"
	"strings"

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

func VariableIncompatibleTypeError(pc ObjectPathContext, varType string, value any) *ValidateError {
	return &ValidateError{
		ObjectPathContext: pc,
		ErrorMessage:      fmt.Sprintf("Specified variable value (%T) is incompatible with expected variable type (%s)", value, varType),
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

func MissingCapabilityNameError(pc ObjectPathContext) *ValidateError {
	return &ValidateError{
		ObjectPathContext: pc,
		ErrorMessage:      fmt.Sprintf("Capability requires a name"),
	}
}

func EnvVariableKeyStartsWithNumberError(pc ObjectPathContext) ValidateError {
	return ValidateError{
		ObjectPathContext: pc,
		ErrorMessage:      "Invalid environment variable, key must not start with a number",
	}
}

func EnvVariableKeyInvalidCharsError(pc ObjectPathContext) ValidateError {
	return ValidateError{
		ObjectPathContext: pc,
		ErrorMessage:      "Invalid environment variable, key must contain only letters, numbers, and underscores",
	}
}

func MissingSubdomainTemplateError(pc ObjectPathContext) ValidateError {
	return ValidateError{
		ObjectPathContext: pc,
		ErrorMessage:      "Subdomain template is required",
	}
}

func InvalidDataClassificationError(pc ObjectPathContext, value string, allowed []string) *ValidateError {
	return &ValidateError{
		ObjectPathContext: pc,
		ErrorMessage:      fmt.Sprintf("Invalid data classification value (%s), must be one of: %s", value, strings.Join(allowed, ", ")),
	}
}

func UnsupportedAppCategoryError(pc ObjectPathContext, moduleSource, subcategory string) ValidateError {
	return ValidateError{
		ObjectPathContext: pc,
		ErrorMessage:      fmt.Sprintf("Module (%s) does not support application category (%s)", moduleSource, subcategory),
	}
}

func InvalidEventActionError(pc ObjectPathContext, actions []string) *ValidateError {
	if len(actions) == 0 {
		return nil
	}
	if len(actions) == 1 {
		return &ValidateError{
			ObjectPathContext: pc,
			ErrorMessage:      fmt.Sprintf("Event Action (%s) is not a valid event action", actions[0]),
		}
	}
	return &ValidateError{
		ObjectPathContext: pc,
		ErrorMessage:      fmt.Sprintf("Event Actions (%s) are not valid event actions", strings.Join(actions, ",")),
	}
}

func InvalidEventStatusError(pc ObjectPathContext, actions []string) *ValidateError {
	if len(actions) == 0 {
		return nil
	}
	if len(actions) == 1 {
		return &ValidateError{
			ObjectPathContext: pc,
			ErrorMessage:      fmt.Sprintf("Event Status (%s) is not a valid event status", actions[0]),
		}
	}
	return &ValidateError{
		ObjectPathContext: pc,
		ErrorMessage:      fmt.Sprintf("Event Statuses (%s) are not valid event statuses", strings.Join(actions, ",")),
	}
}

func InvalidEventTargetError(pc ObjectPathContext, target string) ValidateError {
	return ValidateError{
		ObjectPathContext: pc,
		ErrorMessage:      fmt.Sprintf("Event Target (%s) is not a valid event target", target),
	}
}

func InvalidBuildValueError(pc ObjectPathContext, value string, allowed []string) ValidateError {
	return ValidateError{
		ObjectPathContext: pc,
		ErrorMessage:      fmt.Sprintf("Invalid value (%s), must be one of: %s", value, strings.Join(allowed, ", ")),
	}
}

func InvalidBuildPathError(pc ObjectPathContext, value string) ValidateError {
	return ValidateError{
		ObjectPathContext: pc,
		ErrorMessage:      fmt.Sprintf("Invalid path (%s), must be a relative path inside the repository", value),
	}
}

func BuildSectionNotUsedError(pc ObjectPathContext, section, builder string) ValidateError {
	return ValidateError{
		ObjectPathContext: pc,
		ErrorMessage:      fmt.Sprintf("Settings of the %s builder are not used, the builder is (%s)", section, builder),
	}
}

func BuildSettingNotUsedError(pc ObjectPathContext, packageMode string) ValidateError {
	return ValidateError{
		ObjectPathContext: pc,
		ErrorMessage:      fmt.Sprintf("Setting is not used by the package mode (%s) of the app's module", packageMode),
	}
}

func BuilderCannotBuildError(pc ObjectPathContext, builder, packageMode string) ValidateError {
	return ValidateError{
		ObjectPathContext: pc,
		ErrorMessage:      fmt.Sprintf("The %s builder cannot build the package mode (%s) of the app's module", builder, packageMode),
	}
}

// DockerfileOutsideRepoError reports a dockerfile that leaves the repository once it is resolved against docker_context
func DockerfileOutsideRepoError(pc ObjectPathContext, dockerfile, dockerContext string) ValidateError {
	return ValidateError{
		ObjectPathContext: pc,
		ErrorMessage:      fmt.Sprintf("Invalid path (%s), dockerfile is relative to docker_context (%s) and must stay inside the repository", dockerfile, dockerContext),
	}
}
