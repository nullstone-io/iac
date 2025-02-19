package iac

import (
	"github.com/nullstone-io/iac/core"
)

// Validate performs validation on the parsed IaC files
// This operation only performs partial validation if done before Resolve
func Validate(input ConfigFiles) core.ValidateErrors {
	errs := core.ValidateErrors{}
	if input.Config != nil {
		for _, err := range input.Config.Validate() {
			err.IacContext = input.Config.IacContext
			errs = append(errs, err)
		}
	}

	for _, cur := range input.Overrides {
		for _, err := range cur.Validate() {
			err.IacContext = cur.IacContext
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}
