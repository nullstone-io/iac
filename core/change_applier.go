package core

import (
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

type ChangeApplier interface {
	ApplyChangesTo(ic IacContext, updater WorkspaceConfigUpdater) error
}

type WorkspaceConfigUpdater interface {
	UpdateSchema(moduleSource, moduleConstraint string, moduleVersion *types.ModuleVersion)
	UpdateVariableValue(name string, value any)
	UpdateConnectionTarget(name string, desired, effective types.ConnectionTarget)
	AddOrUpdateEnvVariable(name string, value string, sensitive bool)
	RemoveEnvVariablesNotIn(envVariables map[string]string)
	GetCapabilityUpdater(identity CapabilityIdentity) CapabilityConfigUpdater
	AddCapability(id int64, name string) CapabilityConfigUpdater
	RemoveCapabilitiesNotIn(identities CapabilityIdentities)
	UpdateDomainName(domainNameTemplate *string)
	UpdateSubdomainName(domainNameTemplate, subdomainNameTemplate *string, reservation *types.SubdomainReservation)
	// UpdateDataClassification sets the workspace's data-classification level
	// (empty clears it).
	UpdateDataClassification(level types.ClassificationLevel)
}

// BuildChangeApplier is implemented by configuration that declares `build:` (applications)
type BuildChangeApplier interface {
	ApplyBuildChangesTo(ic IacContext, updater BuildConfigUpdater) error
}

// BuildConfigUpdater receives the build config that IaC files declare for an application.
// It is separate from WorkspaceConfigUpdater because a build config is not part of the workspace config:
// it applies to the next build, with no workspace change and no infra run.
type BuildConfigUpdater interface {
	// ResetBuildConfig returns the build config to the default of the module of the application.
	// config.yml declares the whole build config, so it starts from the default.
	ResetBuildConfig()
	// UpdateBuildConfig changes the build config in place and records ic as the file that declared it
	UpdateBuildConfig(ic IacContext, fn func(config *types.BuildConfig))
}

type ExtraWorkspaceConfigInput struct {
	DnsName    *string
	DomainName *string
	Fqdn       *string
}

type CapabilityConfigUpdater interface {
	UpdateSchema(moduleSource, moduleConstraint string, moduleVersion *types.ModuleVersion)
	UpdateVariableValue(name string, value any)
	UpdateConnectionTarget(name string, desired, effective types.ConnectionTarget)
	UpdateNamespace(namespace *string)
}
