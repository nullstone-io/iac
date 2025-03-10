package config

import (
	"context"
	"github.com/google/go-cmp/cmp"
	"github.com/gorilla/mux"
	"github.com/nullstone-io/iac/core"
	"github.com/nullstone-io/iac/services"
	"github.com/nullstone-io/iac/services/oracle"
	config2 "github.com/nullstone-io/iac/yaml"
	"github.com/nullstone-io/module/config"
	"github.com/stretchr/testify/assert"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
	"os"
	"testing"
)

type FactoryDefaults struct {
	OrgName string
	StackId int64
	BlockId int64
	EnvId   int64
}

func ptr[T any](t T) *T {
	return &t
}

func TestConvertConfiguration(t *testing.T) {
	providerType := "aws"
	defaults := FactoryDefaults{
		OrgName: "nullstone",
		StackId: 123,
		BlockId: 456,
		EnvId:   1,
	}
	latest := "latest"
	primary := "primary"
	subdomainName := "ns-sub-for-acme-docs"
	modules := []*types.Module{
		{
			OrgName:       defaults.OrgName,
			Name:          "aws-fargate-service",
			Category:      "app",
			Subcategory:   "container",
			ProviderTypes: []string{"aws"},
			Platform:      "ecs",
			Subplatform:   "",
			LatestVersion: &types.ModuleVersion{
				Version: "0.0.1",
				Manifest: config.Manifest{
					Connections: map[string]config.Connection{
						"cluster-namespace": {
							Contract: "cluster-namespace/aws/fargate",
							Optional: false,
						},
					},
					Variables: map[string]config.Variable{
						"num_tasks": {
							Type:    "number",
							Default: 1,
						},
						"cpu": {
							Type:    "number",
							Default: 256,
						},
						"memory": {
							Type:    "number",
							Default: 512,
						},
					},
				},
			},
		},
		{
			OrgName:       defaults.OrgName,
			Name:          "aws-s3-site",
			Category:      "app",
			Subcategory:   "static-site",
			ProviderTypes: []string{"aws"},
			Platform:      "s3",
			Subplatform:   "",
			LatestVersion: &types.ModuleVersion{
				Version: "0.0.1",
				Manifest: config.Manifest{
					Connections: nil,
					Variables: map[string]config.Variable{
						"enable_versioned_assets": {
							Type:      "bool",
							Default:   true,
							Sensitive: false,
						},
						"env_vars": {
							Type:      "map(string)",
							Default:   map[string]string{},
							Sensitive: false,
						},
						"env_vars_filename": {
							Type:      "string",
							Default:   "env.json",
							Sensitive: false,
						},
					},
					EnvVariables: nil,
				},
			},
		},
		{
			OrgName:       defaults.OrgName,
			Name:          "aws-load-balancer",
			Category:      "capability",
			ProviderTypes: []string{"aws"},
			AppCategories: []string{"container"},
			LatestVersion: &types.ModuleVersion{
				Version: "0.0.1",
				Manifest: config.Manifest{
					Variables: map[string]config.Variable{
						"enable_https": {
							Type:    "bool",
							Default: true,
						},
						"health_check_enabled": {
							Type:    "bool",
							Default: true,
						},
						"health_check_path": {
							Type:    "string",
							Default: "/",
						},
						"health_check_matcher": {
							Type:    "string",
							Default: "200-499",
						},
						"health_check_healthy_threshold": {
							Type:    "number",
							Default: 2,
						},
						"health_check_unhealthy_threshold": {
							Type:    "number",
							Default: 2,
						},
						"health_check_interval": {
							Type:    "number",
							Default: 5,
						},
						"health_check_timeout": {
							Type:    "number",
							Default: 4,
						},
					},
					Connections: map[string]config.Connection{
						"subdomain": {
							Contract: "subdomain/aws/route53",
						},
					},
				},
			},
		},
		{
			OrgName:       defaults.OrgName,
			Name:          "aws-s3-cdn",
			Category:      "capability",
			ProviderTypes: []string{"aws"},
			AppCategories: []string{"static-site"},
			LatestVersion: &types.ModuleVersion{
				Version: "0.0.1",
				Manifest: config.Manifest{
					Connections: map[string]config.Connection{
						"subdomain": {
							Contract: "subdomain/aws/route53",
						},
					},
					Variables: map[string]config.Variable{
						"enable_www": {
							Type:      "bool",
							Default:   true,
							Sensitive: false,
						},
						"notfound_behavior": {
							Type:      "object({ enabled : bool spa_mode : bool document : string })",
							Default:   map[string]any{"document": "404.html", "enabled": true, "spa_mode": true},
							Sensitive: false,
						},
					},
				},
			},
		},
		{
			OrgName:       defaults.OrgName,
			Name:          "aws-postgres-access",
			Category:      "capability",
			ProviderTypes: []string{"aws"},
			AppCategories: []string{"container", "serverless", "server"},
			LatestVersion: &types.ModuleVersion{
				Version: "0.0.1",
				Manifest: config.Manifest{
					Connections: map[string]config.Connection{
						"postgres": {
							Contract: "datastore/aws/postgres:*",
						},
					},
					Variables: map[string]config.Variable{
						"database_name": {
							Type:      "string",
							Default:   nil,
							Sensitive: false,
						},
					},
				},
			},
		},
		{
			OrgName:       defaults.OrgName,
			Name:          "aws-rds-postgres",
			Category:      "datastore",
			Platform:      "postgres",
			Subplatform:   "rds",
			ProviderTypes: []string{"aws"},
			LatestVersion: &types.ModuleVersion{
				Version: "0.0.1",
				Manifest: config.Manifest{
					Connections: map[string]config.Connection{
						"network": {
							Contract: "network/aws/vpc",
						},
					},
				},
			},
		},
	}
	namespaceBlock := types.Block{
		IdModel:             types.IdModel{Id: 100},
		OrgName:             defaults.OrgName,
		StackId:             defaults.StackId,
		Name:                "namespace0",
		ModuleSource:        "nullstone/aws-fargate-namespace",
		ModuleSourceVersion: latest,
	}
	subdomainBlock := types.Block{
		IdModel:             types.IdModel{Id: 98},
		OrgName:             defaults.OrgName,
		StackId:             defaults.StackId,
		Name:                subdomainName,
		ModuleSource:        "nullstone/aws-autogen-subdomain",
		ModuleSourceVersion: latest,
	}
	postgresBlock := types.Block{
		IdModel:             types.IdModel{Id: 97},
		OrgName:             defaults.OrgName,
		StackId:             defaults.StackId,
		Name:                "postgres",
		ModuleSource:        "nullstone/aws-rds-postgres",
		ModuleSourceVersion: latest,
	}
	blocksById := map[int64]types.Block{namespaceBlock.Id: namespaceBlock, subdomainBlock.Id: subdomainBlock, postgresBlock.Id: postgresBlock}
	blocksByName := map[string]types.Block{namespaceBlock.Name: namespaceBlock, subdomainBlock.Name: subdomainBlock, postgresBlock.Name: postgresBlock}

	tests := []struct {
		name             string
		filename         string
		isOverrides      bool
		want             *EnvConfiguration
		resolveErrors    core.ResolveErrors
		validationErrors core.ValidateErrors
	}{
		{
			name:     "valid configuration",
			filename: "test-fixtures/config.yml",
			want: &EnvConfiguration{
				IacContext: core.IacContext{
					RepoUrl:  "https://github.com/acme/api",
					RepoName: "acme/api",
					Filename: "config.yml",
					Version:  "0.1",
				},
				Events: EventConfigurations{
					"deployments": {
						Name:       "deployments",
						Actions:    []types.EventAction{types.EventActionAppDeployed},
						BlockNames: []string{"acme-docs"},
						Statuses:   []types.EventStatus{types.EventStatusCompleted},
						Targets: EventTargetConfigurations{
							"slack": {
								Target: "slack",
								SlackData: &SlackEventTargetData{
									Channels: []string{"deployments"},
								},
							},
						},
						Blocks: nil,
					},
				},
				Applications: map[string]*AppConfiguration{
					"acme-docs": {
						BlockConfiguration: BlockConfiguration{
							Type:             BlockTypeApplication,
							Category:         types.CategoryApp,
							Name:             "acme-docs",
							ModuleSource:     "nullstone/aws-fargate-service",
							ModuleConstraint: latest,
							Variables: VariableConfigurations{
								"num_tasks": {Value: 2},
							},
							Connections: ConnectionConfigurations{
								"cluster-namespace": {
									DesiredTarget: types.ConnectionTarget{
										BlockName: "namespace0",
									},
								},
							},
						},
						EnvVariables: map[string]string{
							"TESTING": "abc123",
							"BLAH":    "blahblahblah",
						},
						Capabilities: CapabilityConfigurations{
							{
								ModuleSource:     "nullstone/aws-load-balancer",
								ModuleConstraint: latest,
								Variables: VariableConfigurations{
									"health_check_path": {Value: "/status"},
								},
								Connections: ConnectionConfigurations{
									"subdomain": {
										DesiredTarget: types.ConnectionTarget{
											BlockName: subdomainName,
										},
									},
								},
								Namespace: &primary,
							},
						},
					},
				},
				Datastores:        map[string]*DatastoreConfiguration{},
				Subdomains:        map[string]*SubdomainConfiguration{},
				Domains:           map[string]*DomainConfiguration{},
				Ingresses:         map[string]*IngressConfiguration{},
				ClusterNamespaces: map[string]*ClusterNamespaceConfiguration{},
				Clusters:          map[string]*ClusterConfiguration{},
				Networks:          map[string]*NetworkConfiguration{},
				Blocks:            map[string]*BlockConfiguration{},
			},
			resolveErrors:    core.ResolveErrors(nil),
			validationErrors: core.ValidateErrors(nil),
		},
		{
			name:     "app module missing",
			filename: "test-fixtures/config.invalid1.yml",
			want:     nil,
			resolveErrors: core.ResolveErrors{
				{
					ObjectPathContext: core.ObjectPathContext{Path: "apps.acme-docs"},
					ErrorMessage:      "Module is required",
				},
			},
			validationErrors: core.ValidateErrors(nil),
		},
		{
			name:     "invalid app module",
			filename: "test-fixtures/config.invalid2.yml",
			want:     nil,
			resolveErrors: core.ResolveErrors{
				{
					ObjectPathContext: core.ObjectPathContext{Path: "apps.acme-docs", Field: "module"},
					ErrorMessage:      "Module (nullstone/aws-invalid-module) does not exist",
				},
			},
			validationErrors: core.ValidateErrors(nil),
		},
		{
			name:     "not an app module",
			filename: "test-fixtures/config.invalid3.yml",
			want:     nil,
			resolveErrors: core.ResolveErrors{
				{
					ObjectPathContext: core.ObjectPathContext{Path: "apps.acme-docs", Field: "module"},
					ErrorMessage:      "Module (nullstone/aws-s3-cdn) must be app module and match the contract (app/*/*), it is defined as capability/aws/",
				},
			},
			validationErrors: core.ValidateErrors(nil),
		},
		{
			name:          "invalid app variable",
			filename:      "test-fixtures/config.invalid4.yml",
			want:          nil,
			resolveErrors: core.ResolveErrors(nil),
			validationErrors: core.ValidateErrors{
				core.ValidateError{
					ObjectPathContext: core.ObjectPathContext{Path: "apps.acme-docs", Field: "vars", Key: "service_count"},
					ErrorMessage:      "Variable does not exist on the module (nullstone/aws-fargate-service@0.0.1)",
				},
			},
		},
		{
			name:     "capability module missing",
			filename: "test-fixtures/config.invalid5.yml",
			want:     nil,
			resolveErrors: core.ResolveErrors{
				{
					ObjectPathContext: core.ObjectPathContext{Path: "apps.acme-docs", Field: "capabilities", Index: ptr(0)},
					ErrorMessage:      "Module is required",
				},
			},
			validationErrors: core.ValidateErrors(nil),
		},
		{
			name:     "invalid capability module",
			filename: "test-fixtures/config.invalid6.yml",
			want:     nil,
			resolveErrors: core.ResolveErrors{
				{
					ObjectPathContext: core.ObjectPathContext{Path: "apps.acme-docs.capabilities[0]", Field: "module"},
					ErrorMessage:      "Module (nullstone/aws-invalid-module) does not exist",
				},
			},
			validationErrors: core.ValidateErrors(nil),
		},
		{
			name:     "not a capability module",
			filename: "test-fixtures/config.invalid7.yml",
			want:     nil,
			resolveErrors: core.ResolveErrors{
				{
					ObjectPathContext: core.ObjectPathContext{Path: "apps.acme-docs.capabilities[0]", Field: "module"},
					ErrorMessage:      "Module (nullstone/aws-s3-site) must be capability module and match the contract (capability/aws/*), it is defined as app:static-site/aws/s3",
				},
			},
			validationErrors: core.ValidateErrors(nil),
		},
		{
			name:          "capability does not match app subcategory",
			filename:      "test-fixtures/config.invalid8.yml",
			want:          nil,
			resolveErrors: core.ResolveErrors(nil),
			validationErrors: core.ValidateErrors{
				core.ValidateError{
					ObjectPathContext: core.ObjectPathContext{Path: "apps.acme-docs.capabilities[0]", Field: "module"},
					ErrorMessage:      "Module (nullstone/aws-postgres-access) does not support application category (static-site)",
				},
			},
		},
		{
			name:          "invalid capability variable",
			filename:      "test-fixtures/config.invalid9.yml",
			want:          nil,
			resolveErrors: core.ResolveErrors(nil),
			validationErrors: core.ValidateErrors{
				core.ValidateError{
					ObjectPathContext: core.ObjectPathContext{Path: "apps.acme-docs.capabilities[0]", Field: "vars", Key: "database_name"},
					ErrorMessage:      "Variable does not exist on the module (nullstone/aws-load-balancer@latest)",
				},
			},
		},
		{
			name:     "capability block doesn't exist",
			filename: "test-fixtures/config.invalid10.yml",
			want:     nil,
			resolveErrors: core.ResolveErrors{
				{
					ObjectPathContext: core.ObjectPathContext{Path: "apps.acme-docs.capabilities[0]", Field: "connections", Key: "subdomain"},
					ErrorMessage:      "Connection is invalid, block core/ns-sub-for-blah does not exist",
				},
			},
			validationErrors: core.ValidateErrors(nil),
		},
		{
			name:          "block doesn't match contract for capability connection",
			filename:      "test-fixtures/config.invalid11.yml",
			want:          nil,
			resolveErrors: core.ResolveErrors(nil),
			validationErrors: core.ValidateErrors{
				core.ValidateError{
					ObjectPathContext: core.ObjectPathContext{Path: "apps.acme-docs.capabilities[0]", Field: "connections", Key: "subdomain"},
					ErrorMessage:      "Block (postgres) does not match the required contract (subdomain/aws/route53) for the capability connection",
				},
			},
		},
		{
			name:          "blockName is required",
			filename:      "test-fixtures/config.invalid12.yml",
			want:          nil,
			resolveErrors: core.ResolveErrors(nil),
			validationErrors: core.ValidateErrors{
				core.ValidateError{
					ObjectPathContext: core.ObjectPathContext{Path: "apps.acme-docs.capabilities[0]", Field: "connections", Key: "subdomain"},
					ErrorMessage:      "Connection must have a block_name to identify which block it is connected to",
				},
			},
		},
		{
			name:        "valid previews.yml",
			filename:    "test-fixtures/previews.yml",
			isOverrides: true,
			want: &EnvConfiguration{
				IacContext: core.IacContext{
					RepoUrl:     "https://github.com/acme/api",
					RepoName:    "acme/api",
					Filename:    "config.yml",
					IsOverrides: true,
					Version:     "0.1",
				},
				Events: EventConfigurations{
					"deployments": {
						Name:       "deployments",
						Actions:    []types.EventAction{types.EventActionAppDeployed},
						BlockNames: []string{"acme-api"},
						Targets: EventTargetConfigurations{
							"slack": {
								Target: "slack",
								SlackData: &SlackEventTargetData{
									Channels: []string{"deployments"},
								},
							},
						},
						Blocks: nil,
					},
				},
				Applications: map[string]*AppConfiguration{
					"acme-api": {
						BlockConfiguration: BlockConfiguration{
							Type:     BlockTypeApplication,
							Category: types.CategoryApp,
							Name:     "acme-api",
							Variables: VariableConfigurations{
								"enable_versioned_assets": {Value: false},
							},
							Connections: ConnectionConfigurations{},
						},
						EnvVariables: map[string]string{
							"TESTING": "abc123",
							"BLAH":    "blahblahblah",
						},
						Capabilities: CapabilityConfigurations{
							{
								ModuleSource:     "nullstone/aws-s3-cdn",
								ModuleConstraint: "latest",
								Variables:        VariableConfigurations{"enable_www": {Value: true}},
								Namespace:        ptr("secondary"),
								Connections: ConnectionConfigurations{
									"subdomain": {
										DesiredTarget: types.ConnectionTarget{
											StackId:   0,
											StackName: "",
											BlockId:   0,
											BlockName: "ns-sub-for-acme-docs",
											EnvId:     nil,
											EnvName:   "",
										},
									},
								},
							},
						},
					},
				},
				Blocks:            map[string]*BlockConfiguration{},
				ClusterNamespaces: map[string]*ClusterNamespaceConfiguration{},
				Clusters:          map[string]*ClusterConfiguration{},
				Datastores:        map[string]*DatastoreConfiguration{},
				Domains:           map[string]*DomainConfiguration{},
				Ingresses:         map[string]*IngressConfiguration{},
				Networks:          map[string]*NetworkConfiguration{},
				Subdomains:        map[string]*SubdomainConfiguration{},
			},
			resolveErrors:    core.ResolveErrors(nil),
			validationErrors: core.ValidateErrors(nil),
		},
	}

	router := mux.NewRouter()
	oracle.MockGetModuleVersions(router, modules...)
	apiHub := services.MockApiHub(t, router)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buf, err := os.ReadFile(test.filename)
			assert.NoError(t, err)

			parsed, err := config2.ParseEnvConfiguration(buf)
			assert.NoError(t, err)

			got := ConvertConfiguration("https://github.com/acme/api", "acme/api", "config.yml", test.isOverrides, *parsed)

			if test.want != nil {
				if diff := cmp.Diff(test.want, got); diff != "" {
					t.Errorf("(-want, +got): %s", diff)
				}
			}

			sr := &find.StackResolver{
				ApiClient:           apiHub.Client(defaults.OrgName),
				Stack:               types.Stack{Name: "core", ProviderType: providerType},
				PreviewsSharedEnvId: 0,
				EnvsById:            map[int64]types.Environment{},
				EnvsByName:          map[string]types.Environment{},
				BlocksById:          blocksById,
				BlocksByName:        blocksByName,
			}
			resolver := core.NewApiResolver(apiHub.Client(defaults.OrgName), defaults.StackId, defaults.EnvId)
			resolver.ResourceResolver.StacksById[defaults.StackId] = sr
			resolver.ResourceResolver.StacksByName["core"] = sr
			resolver.EventChannelResolver = core.StaticEventChannelResolver{
				ChannelsByTool: map[string][]map[string]any{
					string(types.IntegrationToolSlack): {
						map[string]any{
							"id":   "C01DBR86SRK",
							"name": "deployments",
						},
						map[string]any{
							"id":   "C01DBR86STK",
							"name": "random",
						},
					},
				},
			}

			ctx := context.Background()
			err1 := got.Resolve(ctx, resolver)
			assert.Equal(t, test.resolveErrors, err1)
			err2 := got.Validate()
			assert.Equal(t, test.validationErrors, err2)
		})
	}
}
