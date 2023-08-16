package overrides

import (
	"fmt"
	"github.com/nullstone-io/iac/config"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
)

type CapabilityOverrides []config.CapabilityConfiguration

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

func (s CapabilityOverrides) FindCapability(toFind models.CapabilityConfig) (*config.CapabilityConfiguration, error) {
	for _, iacCap := range s {
		if toFind.Source != iacCap.ModuleSource {
			continue
		}
		// for each of the connections, find a match in this config's connections
		//   if all the connections match, return this config
		//   if we can't find a match, return an error
		if err := iacCap.matchAllConnections(toFind.Connections); err != nil {
			return nil, fmt.Errorf("found a capability in Nullstone for %q, but the connections do not match: %w", iacCap.ModuleSource, err)
		}
		// the loop above will return an error if we can't find a match
		// if we get this far, we have a match
		return &iacCap, nil
	}
	return nil, nil
}
