package config

import (
	"context"
	"testing"

	"github.com/gorilla/mux"
	"github.com/nullstone-io/iac/core"
	"github.com/nullstone-io/iac/services"
	"github.com/nullstone-io/iac/services/oracle"
	"github.com/nullstone-io/module/config"
	"github.com/stretchr/testify/assert"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

// TestFindBlockModuleInIacCrossEnv covers shared blocks whose connections resolve to a different env
// (e.g. `previews-shared` when testing from a preview env): a same-stack block declared in the IaC session
// is authoritative for module identity regardless of the resolved env.
func TestFindBlockModuleInIacCrossEnv(t *testing.T) {
	namespaceModule := &types.Module{
		OrgName:  "nullstone",
		Name:     "aws-fargate-namespace",
		Category: types.CategoryClusterNamespace,
	}
	base := &EnvConfiguration{
		ClusterNamespaces: map[string]*ClusterNamespaceConfiguration{
			"namespace0": {
				BlockConfiguration: BlockConfiguration{
					Name:         "namespace0",
					ModuleSource: "nullstone/aws-fargate-namespace",
					Module:       namespaceModule,
				},
			},
		},
	}
	finder := NewIacFinder(base, nil, 123, 1)

	tests := []struct {
		name   string
		target types.ConnectionTarget
		want   *types.Module
	}{
		{
			name:   "same stack, same env",
			target: types.ConnectionTarget{StackId: 123, BlockName: "namespace0", EnvId: ptr(int64(1))},
			want:   namespaceModule,
		},
		{
			name:   "same stack, resolved to previews-shared env",
			target: types.ConnectionTarget{StackId: 123, BlockName: "namespace0", EnvId: ptr(int64(99))},
			want:   namespaceModule,
		},
		{
			name:   "different stack",
			target: types.ConnectionTarget{StackId: 456, BlockName: "namespace0", EnvId: ptr(int64(1))},
			want:   nil,
		},
		{
			name:   "not in iac session",
			target: types.ConnectionTarget{StackId: 123, BlockName: "unknown", EnvId: ptr(int64(1))},
			want:   nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := finder.FindBlockModuleInIac(context.Background(), test.target)
			assert.Equal(t, test.want, got)
		})
	}
}

// TestResolveConnectionUnprovisionedWorkspace covers connections that resolve to a workspace with no
// desired config yet (never launched). The resolver must fall back to the block record's intended module
// instead of failing; only a block with no module at all is an error.
func TestResolveConnectionUnprovisionedWorkspace(t *testing.T) {
	orgName := "nullstone"
	stackId := int64(123)
	previewEnvId := int64(1)
	sharedEnvId := int64(99)

	namespaceModule := &types.Module{
		OrgName:       orgName,
		Name:          "aws-fargate-namespace",
		Category:      types.CategoryClusterNamespace,
		ProviderTypes: []string{"aws"},
		Platform:      "fargate",
		LatestVersion: &types.ModuleVersion{Version: "0.0.1", Manifest: config.Manifest{}},
	}
	namespaceBlock := types.Block{
		IdModel:      types.IdModel{Id: 100},
		OrgName:      orgName,
		StackId:      stackId,
		Name:         "namespace0",
		IsShared:     true,
		ModuleSource: "nullstone/aws-fargate-namespace",
	}
	moduleLessBlock := types.Block{
		IdModel:  types.IdModel{Id: 101},
		OrgName:  orgName,
		StackId:  stackId,
		Name:     "no-module0",
		IsShared: true,
	}

	router := mux.NewRouter()
	oracle.MockGetModuleVersions(router, namespaceModule)
	// No desired configs exist: every /configs/latest lookup misses (never-launched workspaces)
	services.MockLatestWorkspaceConfigs(router, nil)
	apiHub := services.MockApiHub(t, router)

	previewEnv := types.Environment{IdModel: types.IdModel{Id: previewEnvId}, Type: types.EnvTypePreview, Name: "poc01", OrgName: orgName, StackId: stackId}
	sharedEnv := types.Environment{IdModel: types.IdModel{Id: sharedEnvId}, Type: types.EnvTypePreviewsShared, Name: "previews-shared", OrgName: orgName, StackId: stackId}
	sr := &find.StackResolver{
		ApiClient:           apiHub.Client(orgName),
		Stack:               types.Stack{IdModel: types.IdModel{Id: stackId}, Name: "core", ProviderType: "aws"},
		PreviewsSharedEnvId: sharedEnvId,
		EnvsById:            map[int64]types.Environment{previewEnvId: previewEnv, sharedEnvId: sharedEnv},
		EnvsByName:          map[string]types.Environment{previewEnv.Name: previewEnv, sharedEnv.Name: sharedEnv},
		BlocksById:          map[int64]types.Block{namespaceBlock.Id: namespaceBlock, moduleLessBlock.Id: moduleLessBlock},
		BlocksByName:        map[string]types.Block{namespaceBlock.Name: namespaceBlock, moduleLessBlock.Name: moduleLessBlock},
	}
	resolver := core.NewApiResolver(apiHub.Client(orgName), stackId, previewEnvId)
	resolver.ResourceResolver.StacksById[stackId] = sr
	resolver.ResourceResolver.StacksByName["core"] = sr

	// The connected blocks are not declared in this IaC session
	finder := NewIacFinder(nil, nil, stackId, previewEnvId)
	pc := core.NewObjectPathContextKey("apps", "acme-api")

	t.Run("falls back to block record module", func(t *testing.T) {
		conn := &ConnectionConfiguration{
			DesiredTarget:   types.ConnectionTarget{BlockName: namespaceBlock.Name},
			EffectiveTarget: types.ConnectionTarget{StackId: stackId, StackName: "core", BlockId: namespaceBlock.Id, BlockName: namespaceBlock.Name, EnvId: ptr(sharedEnvId), EnvName: sharedEnv.Name},
			Schema:          &config.Connection{Contract: "cluster-namespace/aws/fargate"},
		}
		err := conn.Resolve(context.Background(), resolver, finder, pc.SubKey("connections", "cluster-namespace"))
		assert.Nil(t, err)
		if assert.NotNil(t, conn.Module) {
			assert.Equal(t, namespaceModule.Name, conn.Module.Name)
		}
	})

	t.Run("block without module is an error", func(t *testing.T) {
		conn := &ConnectionConfiguration{
			DesiredTarget:   types.ConnectionTarget{BlockName: moduleLessBlock.Name},
			EffectiveTarget: types.ConnectionTarget{StackId: stackId, StackName: "core", BlockId: moduleLessBlock.Id, BlockName: moduleLessBlock.Name, EnvId: ptr(sharedEnvId), EnvName: sharedEnv.Name},
			Schema:          &config.Connection{Contract: "cluster-namespace/aws/fargate"},
		}
		err := conn.Resolve(context.Background(), resolver, finder, pc.SubKey("connections", "cluster-namespace"))
		if assert.NotNil(t, err) {
			expected := core.ResolvedBlockMissingModuleError(pc.SubKey("connections", "cluster-namespace"), "core", moduleLessBlock.Name)
			assert.Equal(t, expected.ErrorMessage, err.ErrorMessage)
		}
	})
}
