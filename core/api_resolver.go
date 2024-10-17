package core

import (
	"context"
	"gopkg.in/nullstone-io/go-api-client.v0"
	"gopkg.in/nullstone-io/go-api-client.v0/artifacts"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

var (
	_ ResolveResolver    = &ApiResolver{}
	_ ConnectionResolver = &ApiResolver{}
)

type ApiResolver struct {
	ApiClient        *api.Client
	ResourceResolver *find.ResourceResolver
}

func NewApiResolver(apiClient *api.Client, stackId, envId int64) *ApiResolver {
	return &ApiResolver{
		ApiClient: apiClient,
		ResourceResolver: &find.ResourceResolver{
			ApiClient:    apiClient,
			CurStackId:   stackId,
			CurEnvId:     envId,
			StacksById:   map[int64]*find.StackResolver{},
			StacksByName: map[string]*find.StackResolver{},
		},
	}
}

func (a *ApiResolver) ResolveBlock(ctx context.Context, ct types.ConnectionTarget) (types.Block, error) {
	return a.ResourceResolver.FindBlock(ctx, ct)
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
