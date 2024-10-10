package iac

import (
	"context"
	"github.com/nullstone-io/iac/core"
)

func Resolve(ctx context.Context, input ParseMapResult, resolver core.ModuleVersionResolver) core.ResolveErrors {
	errs := core.ResolveErrors{}
	if input.Config != nil {
		errs = append(errs, input.Config.Resolve(ctx, resolver)...)
	}
	for _, cur := range input.Overrides {
		errs = append(errs, cur.Resolve(ctx, resolver)...)
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}
