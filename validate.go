package iac

import (
	"context"
	"github.com/BSick7/go-api/errors"
	"github.com/nullstone-io/iac/core"
)

func Validate(ctx context.Context, input ParseMapResult, resolver core.ValidateResolver) errors.ValidationErrors {
	ve := errors.ValidationErrors{}
	if input.Config != nil {
		ve = append(ve, input.Config.Validate(ctx, resolver)...)
	}
	for _, envOverrides := range input.Overrides {
		ve = append(ve, envOverrides.Validate(ctx, resolver)...)
	}
	if len(ve) > 0 {
		return ve
	}
	return nil
}
