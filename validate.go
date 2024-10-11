package iac

import (
	"context"
	"github.com/nullstone-io/iac/core"
)

func Validate(ctx context.Context, input ParseMapResult, resolver core.ValidateResolver) core.ValidateErrors {
	errs := core.ValidateErrors{}
	if input.Config != nil {
		for _, err := range input.Config.Validate(ctx, resolver) {
			err.IacContext = input.Config.IacContext
			errs = append(errs, err)
		}
	}

	for _, cur := range input.Overrides {
		for _, err := range cur.Validate(ctx, resolver) {
			err.IacContext = cur.IacContext
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}
