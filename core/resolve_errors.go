package core

import (
	"fmt"
	"github.com/BSick7/go-api/errors"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

var (
	_ error = ResolveError{}
)

type ResolveError struct {
	IacContext        IacContext        `json:"iacContext"`
	ObjectPathContext ObjectPathContext `json:"objectPathContext"`
	ErrorMessage      string            `json:"errorMessage"`
}

func (e ResolveError) Error() string {
	return fmt.Sprintf("%s => %s", e.IacContext.Context(e.ObjectPathContext), e.ErrorMessage)
}

func (e ResolveError) ToValidationError() errors.ValidationError {
	return errors.ValidationError{
		Context: e.IacContext.Context(e.ObjectPathContext),
		Message: e.ErrorMessage,
	}
}

type ResolveErrors []ResolveError

func (s ResolveErrors) ToValidationErrors() errors.ValidationErrors {
	if len(s) == 0 {
		return nil
	}
	ve := errors.ValidationErrors{}
	for _, re := range s {
		ve = append(ve, re.ToValidationError())
	}
	return ve
}

func InvalidResolveModuleFormatError(pc ObjectPathContext, moduleSource string) *ResolveError {
	return &ResolveError{
		ObjectPathContext: pc.SubField("module"),
		ErrorMessage:      fmt.Sprintf("Invalid module format (%s) - must be in the format \"<module-org>/<module-name>\"", moduleSource),
	}
}

func RequiredModuleError(pc ObjectPathContext) *ResolveError {
	return &ResolveError{
		ObjectPathContext: pc,
		ErrorMessage:      "Module is required",
	}
}

func ModuleVersionLookupFailedError(pc ObjectPathContext, source, version string, err error) *ResolveError {
	return &ResolveError{
		ObjectPathContext: pc.SubField("module_version"),
		ErrorMessage:      fmt.Sprintf("Module version (%s@%s) lookup failed: %s", source, version, err),
	}
}

func MissingModuleError(pc ObjectPathContext, moduleSource string) *ResolveError {
	return &ResolveError{
		ObjectPathContext: pc.SubField("module"),
		ErrorMessage:      fmt.Sprintf("Module (%s) does not exist", moduleSource),
	}
}

func InvalidModuleContractError(pc ObjectPathContext, moduleSource string, want, got types.ModuleContractName) *ResolveError {
	return &ResolveError{
		ObjectPathContext: pc.SubField("module"),
		ErrorMessage:      fmt.Sprintf("Module (%s) must be %s module and match the contract (%s), it is defined as %s", moduleSource, want.Category, want, got),
	}
}

func MissingModuleVersionError(pc ObjectPathContext, source, version string) *ResolveError {
	return &ResolveError{
		ObjectPathContext: pc.SubField("module_version"),
		ErrorMessage:      fmt.Sprintf("Module version (%s@%s) does not exist", source, version),
	}
}

func MissingRequiredConnectionError(pc ObjectPathContext, connName string) ResolveError {
	return ResolveError{
		ObjectPathContext: pc,
		ErrorMessage:      fmt.Sprintf("Connection (%s) is required", connName),
	}
}

func MissingConnectionTargetError(pc ObjectPathContext, err error) *ResolveError {
	return &ResolveError{
		ObjectPathContext: pc,
		ErrorMessage:      fmt.Sprintf("Connection is invalid, %s", err),
	}
}

func LookupConnectionTargetFailedError(pc ObjectPathContext, err error) *ResolveError {
	return &ResolveError{
		ObjectPathContext: pc,
		ErrorMessage:      fmt.Sprintf("Failed to validate connection, error when looking up connection target: %s", err),
	}
}

func ResolvedBlockMissingModuleError(pc ObjectPathContext, stackName, blockName string) *ResolveError {
	return &ResolveError{
		ObjectPathContext: pc.SubField("module"),
		ErrorMessage:      fmt.Sprintf("Module is required on block (%s/%s) that was resolved by a connection.", stackName, blockName),
	}
}

func InvalidModuleFormatError(pc ObjectPathContext, moduleSource string) *ResolveError {
	return &ResolveError{
		ObjectPathContext: pc.SubField("module"),
		ErrorMessage:      fmt.Sprintf("Invalid module format (%s) - must be in the format \"<module-org>/<module-name>\"", moduleSource),
	}
}

func ModuleLookupFailedError(pc ObjectPathContext, moduleSource string, err error) *ResolveError {
	return &ResolveError{
		ObjectPathContext: pc.SubField("module"),
		ErrorMessage:      fmt.Sprintf("Module (%s) lookup failed: %s", moduleSource, err),
	}
}

func ToolChannelLookupFailedError(pc ObjectPathContext, tool string, err error) ResolveError {
	return ResolveError{
		ObjectPathContext: pc,
		ErrorMessage:      fmt.Sprintf("Failed to look up %s channels: %s", tool, err),
	}
}

func InvalidRandomSubdomainTemplateError(pc ObjectPathContext, template string) ResolveError {
	return ResolveError{
		ObjectPathContext: pc,
		ErrorMessage:      fmt.Sprintf("Invalid subdomain template %q: cannot have specify additional text when using `{{ random() }}`.", template),
	}
}

func FailedSubdomainReservationError(pc ObjectPathContext, requested string, err error) ResolveError {
	return ResolveError{
		ObjectPathContext: pc,
		ErrorMessage:      fmt.Sprintf("Failed to reserve subdomain %q: %s", requested, err),
	}
}
