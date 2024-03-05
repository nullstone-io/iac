package config

import (
	"context"
	"github.com/nullstone-io/iac/core"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

type CapabilityConfigurations []CapabilityConfiguration

func (c CapabilityConfigurations) Normalize(ctx context.Context, resolver *find.ResourceResolver) error {
	for i, iacCap := range c {
		resolved, err := iacCap.Normalize(ctx, resolver)
		if err != nil {
			return err
		}
		c[i] = resolved
	}
	return nil
}

type CapabilityConfiguration struct {
	ModuleSource        string
	ModuleSourceVersion string
	Variables           map[string]any
	Connections         types.ConnectionTargets
	Namespace           *string
}

func (c CapabilityConfiguration) Normalize(ctx context.Context, resolver *find.ResourceResolver) (CapabilityConfiguration, error) {
	if err := core.NormalizeConnectionTargets(ctx, c.Connections, resolver); err != nil {
		return c, err
	}
	return c, nil
}
