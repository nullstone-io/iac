package core

import (
	errors3 "errors"
	"github.com/BSick7/go-api/errors"
)

func AsValidationErrors(err error) errors.ValidationErrors {
	var verrs errors.ValidationErrors
	if errors3.As(err, &verrs) {
		return verrs
	}
	var verr errors.ValidationError
	if errors3.As(err, &verr) {
		return errors.ValidationErrors{verr}
	}
	return nil
}

func appendValidationErrors(ve *errors.ValidationErrors, err error) error {
	if err == nil {
		return nil
	}

	if verrs := AsValidationErrors(err); verrs != nil {
		*ve = append(*ve, verrs...)
	}

	return err
}
