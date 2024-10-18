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
