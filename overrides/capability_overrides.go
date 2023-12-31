package overrides

import (
	"github.com/nullstone-io/iac/core"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
)

type CapabilityOverrides []core.CapabilityConfiguration

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
