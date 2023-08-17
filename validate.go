package iac

import (
	"github.com/BSick7/go-api/errors"
	"github.com/nullstone-io/iac/config"
	"github.com/nullstone-io/iac/overrides"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
)

func Validate(config *config.EnvConfiguration, overrides *overrides.ConfigurationOverrides, resolver *find.ResourceResolver) error {
	ve := errors.ValidationErrors{}
	if config != nil {
		verrs, err := config.Validate(resolver)
		if err != nil {
			return err
		}
		ve = append(ve, verrs...)
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
