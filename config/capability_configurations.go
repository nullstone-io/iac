package config

import (
	"gopkg.in/nullstone-io/go-api-client.v0/find"
)

type CapabilityConfigurations []CapabilityConfiguration

func (s CapabilityConfigurations) Normalize(resolver *find.ResourceResolver) error {
	for i, iacCap := range s {
		resolved, err := iacCap.Normalize(resolver)
		if err != nil {
			return err
		}
		s[i] = resolved
	}
	return nil
}
