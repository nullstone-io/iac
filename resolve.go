package iac

import (
	"context"
	"github.com/nullstone-io/iac/core"
)

func Resolve(ctx context.Context, input ConfigFiles, resolver core.ResolveResolver) core.ResolveErrors {
	errs := core.ResolveErrors{}
	if input.Config != nil {
		for _, err := range input.Config.Resolve(ctx, resolver) {
			err.IacContext = input.Config.IacContext
			errs = append(errs, err)
		}
	}
	for _, cur := range input.Overrides {
		for _, err := range cur.Resolve(ctx, resolver) {
			err.IacContext = cur.IacContext
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}
