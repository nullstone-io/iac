package iac

import (
	"github.com/nullstone-io/iac/config"
	"github.com/nullstone-io/iac/overrides"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
)

func Normalize(config *config.EnvConfiguration, overrides map[string]overrides.EnvOverrides, resolver *find.ResourceResolver) error {
	if config != nil {
		if err := config.Normalize(resolver); err != nil {
			return err
		}
	}

	for _, envOverrides := range overrides {
		if err := envOverrides.Normalize(resolver); err != nil {
			return err
		}
	}

	return nil
}
