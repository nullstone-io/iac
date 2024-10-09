package config

import (
	"context"
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
	ModuleSource        string                  `json:"moduleSource"`
	ModuleSourceVersion string                  `json:"moduleSourceVersion"`
	Variables           map[string]any          `json:"vars"`
	Connections         types.ConnectionTargets `json:"connections"`
	Namespace           *string                 `json:"namespace"`
}

func (c CapabilityConfiguration) Normalize(ctx context.Context, resolver *find.ResourceResolver) (CapabilityConfiguration, error) {
	if err := NormalizeConnectionTargets(ctx, c.Connections, resolver); err != nil {
		return c, err
	}
	return c, nil
}
