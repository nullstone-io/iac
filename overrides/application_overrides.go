package overrides

import (
	"github.com/BSick7/go-api/errors"
	"github.com/nullstone-io/iac/core"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
)

type ApplicationOverrides struct {
	BlockOverrides
	EnvVariables map[string]string   `yaml:"environment"`
	Capabilities CapabilityOverrides `yaml:"capabilities"`
}

func (a *ApplicationOverrides) Validate(resolver *find.ResourceResolver) (errors.ValidationErrors, error) {
	// TODO: Implement: How do we validate if we don't have a module to resolve
	return errors.ValidationErrors{}, nil
}

func (a *ApplicationOverrides) Normalize(resolver *find.ResourceResolver) error {
	err := core.NormalizeConnectionTargets(a.Connections, resolver)
	if err != nil {
		return err
	}
	return a.Capabilities.Normalize(resolver)
}
