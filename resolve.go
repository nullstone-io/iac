package iac

import (
	"context"

	"github.com/nullstone-io/iac/core"
)

// Resolve performs resolution on all blocks and capabilities:
// - Connection targets and verify connection contract
// - Capability connection target block and verify connection contract
func Resolve(ctx context.Context, input ConfigFiles, resolver core.ResolveResolver, iacFinder core.IacFinder) core.ResolveErrors {
	errs := core.ResolveErrors{}
	if input.Config != nil {
		for _, err := range input.Config.Resolve(ctx, resolver, iacFinder) {
			err.IacContext = input.Config.IacContext
			errs = append(errs, err)
		}
	}
	for _, cur := range input.Overrides {
		for _, err := range cur.Resolve(ctx, resolver, iacFinder) {
			err.IacContext = cur.IacContext
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}
