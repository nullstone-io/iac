package iac

import (
	"context"
	"github.com/nullstone-io/iac/config"
	"github.com/nullstone-io/iac/core"
)

func Normalize(ctx context.Context, config *config.EnvConfiguration, overrides map[string]config.EnvConfiguration, resolver core.ConnectionResolver) error {
	if config != nil {
		if err := config.Normalize(ctx, resolver); err != nil {
			return err
		}
	}

	for _, envOverrides := range overrides {
		if err := envOverrides.Normalize(ctx, resolver); err != nil {
			return err
		}
	}

	return nil
}
