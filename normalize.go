package iac

import (
	"context"
	"github.com/nullstone-io/iac/core"
)

func Normalize(ctx context.Context, input ParseMapResult, resolver core.ConnectionResolver) core.NormalizeErrors {
	errs := core.NormalizeErrors{}

	if input.Config != nil {
		for _, err := range input.Config.Normalize(ctx, resolver) {
			err.IacContext = input.Config.IacContext
			errs = append(errs, err)
		}
	}

	for _, cur := range input.Overrides {
		for _, err := range cur.Normalize(ctx, resolver) {
			err.IacContext = input.Config.IacContext
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}
