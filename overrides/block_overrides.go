package overrides

import (
	"github.com/BSick7/go-api/errors"
	"github.com/nullstone-io/iac/core"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
)

type BlockOverrides struct {
	Name        string                 `yaml:"-"`
	Variables   map[string]any         `yaml:"vars"`
	Connections core.ConnectionTargets `yaml:"connections"`
}

func (b *BlockOverrides) Validate(resolver *find.ResourceResolver) (errors.ValidationErrors, error) {
	// TODO: Implement: How do we validate if we don't have a module to resolve
	return errors.ValidationErrors{}, nil
}

func (b *BlockOverrides) Normalize(resolver *find.ResourceResolver) error {
	return core.NormalizeConnectionTargets(b.Connections, resolver)
}
