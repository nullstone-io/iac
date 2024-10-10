package iac

import (
	"context"
	"github.com/nullstone-io/iac/core"
)

func Normalize(ctx context.Context, input ParseMapResult, resolver core.ConnectionResolver) error {
	if input.Config != nil {
		if err := input.Config.Normalize(ctx, resolver); err != nil {
			return err
		}
	}

	for _, envOverrides := range input.Overrides {
		if err := envOverrides.Normalize(ctx, resolver); err != nil {
			return err
		}
	}

	return nil
}
