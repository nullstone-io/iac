package overrides

import (
	"github.com/nullstone-io/iac/core"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

type CapabilityOverrides []CapabilityOverride

func (s CapabilityOverrides) Normalize(resolver *find.ResourceResolver) error {
	for i, iacCap := range s {
		resolved, err := iacCap.Normalize(resolver)
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

func (c CapabilityOverride) Normalize(resolver *find.ResourceResolver) (CapabilityOverride, error) {
	if err := core.NormalizeConnectionTargets(c.Connections, resolver); err != nil {
		return c, err
	}
	return c, nil
}
