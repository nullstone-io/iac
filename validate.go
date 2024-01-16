package iac

import (
	errs "errors"
	"github.com/BSick7/go-api/errors"
	"github.com/nullstone-io/iac/config"
	"github.com/nullstone-io/iac/overrides"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
)

func Validate(config *config.EnvConfiguration, overrides *overrides.EnvOverrides, resolver *find.ResourceResolver) error {
	ve := errors.ValidationErrors{}
	if config != nil {
		err := config.Validate(resolver)
		if err != nil {
			var verrs errors.ValidationErrors
			if errs.As(err, &verrs) {
				ve = append(ve, verrs...)
			} else {
				return err
			}
		}
	}
	if overrides != nil {
		verrs, err := overrides.Validate(resolver)
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
