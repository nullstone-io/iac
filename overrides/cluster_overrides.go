package overrides

import (
	"github.com/BSick7/go-api/errors"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
)

type ClusterOverrides struct {
	Name      string         `yaml:"-"`
	Variables map[string]any `yaml:"vars"`
}

func (c *ClusterOverrides) Validate(resolver *find.ResourceResolver) (errors.ValidationErrors, error) {
	// TODO: Implement: How do we validate if we don't have a module to resolve
	return errors.ValidationErrors{}, nil
}

func (c *ClusterOverrides) Normalize(resolver *find.ResourceResolver) error {
	return nil
}
