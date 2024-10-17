package core

import (
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

type ChangeApplier interface {
	ApplyChangesTo(ic IacContext, updater WorkspaceConfigUpdater) error
}

type WorkspaceConfigUpdater interface {
	UpdateSchema(moduleSource string, moduleVersion *types.ModuleVersion)
	UpdateVariableValue(name string, value any)
	UpdateConnectionTarget(name string, value types.ConnectionTarget)
	AddOrUpdateEnvVariable(name string, value string, sensitive bool)
	RemoveEnvVariablesNotIn(envVariables map[string]string)
	GetCapabilityUpdater(identity CapabilityIdentity) CapabilityConfigUpdater
	RemoveCapabilitiesNotIn(identities CapabilityIdentities)
}

type CapabilityConfigUpdater interface {
	UpdateSchema(moduleSource string, moduleVersion *types.ModuleVersion)
	UpdateVariableValue(name string, value any)
	UpdateConnectionTarget(name string, value types.ConnectionTarget)
	UpdateNamespace(namespace *string)
}

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
	ModuleSource      string                  `json:"moduleSource"`
	ConnectionTargets types.ConnectionTargets `json:"connectionTargets"`
}

func (i CapabilityIdentity) Match(other CapabilityIdentity) bool {
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
