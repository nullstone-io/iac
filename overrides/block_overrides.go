package overrides

import (
	"context"
	"github.com/BSick7/go-api/errors"
	"github.com/nullstone-io/iac/yaml"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
)

type BlockOverrides struct {
	Name      string         `json:"name"`
	Variables map[string]any `json:"variables"`
}

func convertBlockOverrides(parsed map[string]yaml.BlockOverrides) map[string]BlockOverrides {
	result := make(map[string]BlockOverrides)
	for blockName, blockValue := range parsed {
		block := BlockOverrides{
			Name:      blockName,
			Variables: blockValue.Variables,
		}
		result[blockName] = block
	}
	return result
}

func (b *BlockOverrides) Validate(resolver *find.ResourceResolver) (errors.ValidationErrors, error) {
	// TODO: Implement: How do we validate if we don't have a module to resolve
	return errors.ValidationErrors{}, nil
}

func (b *BlockOverrides) Normalize(ctx context.Context, resolver *find.ResourceResolver) error {
	return nil
}
