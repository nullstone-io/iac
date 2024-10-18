package core

import "gopkg.in/nullstone-io/go-api-client.v0/types"

type CapabilityIdentities []CapabilityIdentity

func (s CapabilityIdentities) Find(match CapabilityIdentity) *CapabilityIdentity {
	for _, cur := range s {
		if cur.Match(match) {
			return &cur
		}
	}
	return nil
}

type CapabilityIdentity struct {
	// Name is unique to the application that holds this capability
	Name string `json:"name"`

	// ModuleSource
	// Deprecated
	ModuleSource string `json:"moduleSource"`
	// ConnectionTargets
	// Deprecated
	ConnectionTargets types.ConnectionTargets `json:"connectionTargets"`
}

func (i CapabilityIdentity) Match(other CapabilityIdentity) bool {
	if i.Name != "" && other.Name != "" {
		return i.Name == other.Name
	}

	if i.ModuleSource != other.ModuleSource {
		return false
	}

	if len(i.ConnectionTargets) != len(other.ConnectionTargets) {
		return false
	}
	for name, target := range i.ConnectionTargets {
		otherConn, ok := other.ConnectionTargets[name]
		if !ok {
			return false
		}
		if !otherConn.Match(target) {
			return false
		}
	}
	return true
}
