package core

import (
	"github.com/nullstone-io/iac/config"
	"github.com/nullstone-io/iac/overrides"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
)

func Normalize(config *config.EnvConfiguration, overrides *overrides.ConfigurationOverrides, resolver *find.ResourceResolver) error {
	if config != nil {
		if err := config.Normalize(resolver); err != nil {
			return err
		}
	}
	if overrides != nil {
		if err := overrides.Normalize(resolver); err != nil {
			return err
		}
	}
	return nil
}
