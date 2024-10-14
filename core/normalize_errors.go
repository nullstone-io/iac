package core

import (
	"fmt"
	"github.com/BSick7/go-api/errors"
)

var (
	_ error = NormalizeError{}
)

type NormalizeError struct {
	IacContext        IacContext        `json:"iacContext"`
	ObjectPathContext ObjectPathContext `json:"objectPathContext"`
	ErrorMessage      string            `json:"errorMessage"`
}

func (e NormalizeError) Error() string {
	return fmt.Sprintf("%s => %s", e.IacContext.Context(e.ObjectPathContext), e.ErrorMessage)
}

func (e NormalizeError) ToValidationError() errors.ValidationError {
	return errors.ValidationError{
		Context: e.IacContext.Context(e.ObjectPathContext),
		Message: e.ErrorMessage,
	}
}

type NormalizeErrors []NormalizeError

func (s NormalizeErrors) ToValidationErrors() errors.ValidationErrors {
	if len(s) == 0 {
		return nil
	}
	ve := errors.ValidationErrors{}
	for _, ne := range s {
		ve = append(ve, ne.ToValidationError())
	}
	return ve
}
