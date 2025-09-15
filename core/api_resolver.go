package core

import (
	"context"
	"errors"
	"fmt"

	"gopkg.in/nullstone-io/go-api-client.v0"
	"gopkg.in/nullstone-io/go-api-client.v0/artifacts"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

var (
	_ ResolveResolver    = &ApiResolver{}
	_ ConnectionResolver = &ApiResolver{}

	ErrEventChannelsNotInitialized = errors.New("event channels have not been initialized")
)

type ApiResolver struct {
	ApiClient            *api.Client
	ResourceResolver     *find.ResourceResolver
	EventChannelResolver EventChannelResolver
	// IacFinder helps lookups for resources contained within the IaC file
	IacFinder IacFinder
}

func NewApiResolver(apiClient *api.Client, iacResolver IacFinder, stackId, envId int64) *ApiResolver {
	return &ApiResolver{
		ApiClient: apiClient,
		ResourceResolver: &find.ResourceResolver{
			ApiClient:    apiClient,
			CurStackId:   stackId,
			CurEnvId:     envId,
			StacksById:   map[int64]*find.StackResolver{},
			StacksByName: map[string]*find.StackResolver{},
		},
		EventChannelResolver: &ApiEventChannelResolver{ApiClient: apiClient},
		IacFinder:            iacResolver,
	}
}

func (a *ApiResolver) ResolveBlock(ctx context.Context, ct types.ConnectionTarget) (types.Block, error) {
	return a.ResourceResolver.FindBlock(ctx, ct)
}

func (a *ApiResolver) FindBlockModuleInIac(ctx context.Context, ct types.ConnectionTarget) *types.Module {
	return a.IacFinder.FindBlockModuleInIac(ctx, ct)
}

func (a *ApiResolver) ResolveWorkspaceModuleConfig(ctx context.Context, ct types.ConnectionTarget) (WorkspaceModuleConfig, error) {
	effective, err := a.ResourceResolver.Resolve(ctx, ct)
	if err != nil {
		return WorkspaceModuleConfig{}, err
	}
	wc, err := a.ApiClient.WorkspaceConfigs().GetLatest(ctx, effective.StackId, effective.BlockId, *effective.EnvId)
	if err != nil {
		return WorkspaceModuleConfig{}, err
	} else if wc == nil {
		return WorkspaceModuleConfig{}, nil
	}
	return WorkspaceModuleConfig{
		Module:           wc.Source,
		ModuleConstraint: wc.SourceConstraint,
	}, nil
}

func (a *ApiResolver) ResolveConnection(ctx context.Context, ct types.ConnectionTarget) (types.ConnectionTarget, error) {
	return a.ResourceResolver.Resolve(ctx, ct)
}

func (a *ApiResolver) ResolveModule(ctx context.Context, source artifacts.ModuleSource) (*types.Module, error) {
	return a.ApiClient.Modules().Get(ctx, source.OrgName, source.ModuleName)
}

func (a *ApiResolver) ResolveModuleVersion(ctx context.Context, source artifacts.ModuleSource, version string) (*types.Module, *types.ModuleVersion, error) {
	m, err := a.ResolveModule(ctx, source)
	if err != nil {
		return nil, nil, err
	} else if m == nil {
		return nil, nil, nil
	}

	if version == "latest" {
		return m, m.LatestVersion, nil
	}
	mv, err := a.ApiClient.ModuleVersions().Get(ctx, source.OrgName, source.ModuleName, version)
	return m, mv, err
}

func (a *ApiResolver) ListChannels(ctx context.Context, tool string) ([]map[string]any, error) {
	if a.EventChannelResolver == nil {
		return nil, ErrEventChannelsNotInitialized
	}
	return a.EventChannelResolver.ListChannels(ctx, tool)
}

func (a *ApiResolver) ReserveNullstoneSubdomain(ctx context.Context, blockName string, requested string) (*types.SubdomainReservation, error) {
	block, err := a.ResolveBlock(ctx, types.ConnectionTarget{BlockName: blockName})
	if err != nil {
		return nil, err
	}
	envId := a.ResourceResolver.CurEnvId
	reservation, err := a.ApiClient.Subdomains().ReserveNullstone(ctx, block.StackId, block.Id, envId, requested)
	if err != nil {
		return nil, err
	}
	if reservation == nil {
		return nil, fmt.Errorf("nullstone reservation returned no result")
	}
	return reservation, nil
}
