package core

import (
	errors3 "errors"
	"github.com/BSick7/go-api/errors"
)

// AsValidationErrors is a helper function to format an error into validation errors for the user to see
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
