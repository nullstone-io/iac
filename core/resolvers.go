package core

import (
	"context"

	"gopkg.in/nullstone-io/go-api-client.v0/artifacts"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

type InitializeResolver interface {
	ModuleVersionResolver
}

type ResolveResolver interface {
	BlockResolver
	IacFinder
	WorkspaceModuleConfigResolver
	ModuleResolver
	ModuleVersionResolver
	EventChannelResolver
	SubdomainReserver
}

type NormalizeResolver interface {
	BlockResolver
	ConnectionResolver
}

type BlockResolver interface {
	// ResolveBlock loads the types.Block based on the resolved connection target
	// This relies on the StackId/BlockId and doesn't consider EnvId
	ResolveBlock(ctx context.Context, ct types.ConnectionTarget) (types.Block, error)
}

type IacFinder interface {
	// FindBlockModuleInIac looks for a BlockConfiguration in the iac configuration files for this IacFinder
	// This expects the input connection target to be fully resolved (all fields populated)
	FindBlockModuleInIac(ctx context.Context, ct types.ConnectionTarget) *types.Module
}

type ConnectionResolver interface {
	// ResolveConnection resolves and backfills any missing fields from input types.ConnectionTarget
	ResolveConnection(ctx context.Context, ct types.ConnectionTarget) (types.ConnectionTarget, error)
}

type WorkspaceModuleConfig struct {
	Module           string
	ModuleConstraint string
}

type WorkspaceModuleConfigResolver interface {
	ResolveWorkspaceModuleConfig(ctx context.Context, ct types.ConnectionTarget) (WorkspaceModuleConfig, error)
}

type SubdomainReserver interface {
	ReserveNullstoneSubdomain(ctx context.Context, blockName string, requested string) (*types.SubdomainReservation, error)
}

type ModuleResolver interface {
	// ResolveModule pulls module metadata for the input source
	ResolveModule(ctx context.Context, source artifacts.ModuleSource) (*types.Module, error)
}

type ModuleVersionResolver interface {
	// ResolveModuleVersion pulls module metadata and version metadata for the input source/version
	ResolveModuleVersion(ctx context.Context, source artifacts.ModuleSource, version string) (*types.Module, *types.ModuleVersion, error)
}

type EventChannelResolver interface {
	ListChannels(ctx context.Context, tool string) ([]map[string]any, error)
}
