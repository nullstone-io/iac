package iac

import (
	"context"

	"github.com/nullstone-io/iac/core"
)

// Initialize performs initialization on all blocks and capabilities:
// - Module schema (variables + connections)
// - Capability module schema (variables + connections)
// Initialize is useful as a separate step from Resolve because we need to initialize all Block information before resolving connections
func Initialize(ctx context.Context, input ConfigFiles, resolver core.InitializeResolver) core.InitializeErrors {
	errs := core.InitializeErrors{}
	if input.Config != nil {
		for _, err := range input.Config.Initialize(ctx, resolver) {
			err.IacContext = input.Config.IacContext
			errs = append(errs, err)
		}
	}
	for _, cur := range input.Overrides {
		for _, err := range cur.Initialize(ctx, resolver) {
			err.IacContext = cur.IacContext
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}
