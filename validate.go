package iac

import (
	"context"
	"github.com/BSick7/go-api/errors"
	"github.com/nullstone-io/iac/config"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
)

func Validate(ctx context.Context, config *config.EnvConfiguration, overrides map[string]config.EnvConfiguration, resolver *find.ResourceResolver) errors.ValidationErrors {
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
