package iac

import (
	"context"
	errs "errors"
	"github.com/BSick7/go-api/errors"
	"github.com/nullstone-io/iac/config"
	"github.com/nullstone-io/iac/overrides"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
)

func Validate(ctx context.Context, config *config.EnvConfiguration, overrides map[string]overrides.EnvOverrides, resolver *find.ResourceResolver) error {
	ve := errors.ValidationErrors{}
	if config != nil {
		err := config.Validate(ctx, resolver)
		if err != nil {
			var verrs errors.ValidationErrors
			if errs.As(err, &verrs) {
				ve = append(ve, verrs...)
			} else {
				return err
			}
		}
	}
	for _, envOverrides := range overrides {
		verrs, err := envOverrides.Validate(resolver)
		if err != nil {
			return err
		}
		ve = append(ve, verrs...)
	}

	if len(ve) > 0 {
		return ve
	}
	return nil
}
