package iac

import (
	"context"
	"github.com/BSick7/go-api/errors"
	"github.com/nullstone-io/iac/config"
	"github.com/nullstone-io/iac/core"
)

func Validate(ctx context.Context, config *config.EnvConfiguration, overrides map[string]config.EnvConfiguration, resolver core.ValidateResolver) errors.ValidationErrors {
	ve := errors.ValidationErrors{}
	if config != nil {
		ve = append(ve, config.Validate(ctx, resolver)...)
	}
	for _, envOverrides := range overrides {
		ve = append(ve, envOverrides.Validate(ctx, resolver)...)
	}
	if len(ve) > 0 {
		return ve
	}
	return nil
}
