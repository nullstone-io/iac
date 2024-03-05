package overrides

import (
	"context"
	"github.com/nullstone-io/iac/core"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

type CapabilityOverrides []CapabilityOverride

func (s CapabilityOverrides) Normalize(ctx context.Context, resolver *find.ResourceResolver) error {
	for i, iacCap := range s {
		resolved, err := iacCap.Normalize(ctx, resolver)
		if err != nil {
			return err
		}
		s[i] = resolved
	}
	return nil
}

type CapabilityOverride struct {
	ModuleSource        string
	ModuleSourceVersion string
	Variables           map[string]any
	Connections         types.ConnectionTargets
	Namespace           *string
}

func (c CapabilityOverride) Normalize(ctx context.Context, resolver *find.ResourceResolver) (CapabilityOverride, error) {
	if err := core.NormalizeConnectionTargets(ctx, c.Connections, resolver); err != nil {
		return c, err
	}
	return c, nil
}
