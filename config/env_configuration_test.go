package config

import (
	"github.com/BSick7/go-api/errors"
	"github.com/gorilla/mux"
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
		want             *EnvConfiguration
		validationErrors error
	}{
		{
			name:     "valid configuration",
			filename: "test-fixtures/config.yml",
			want: &EnvConfiguration{
				RepoName: "acme/api",
				Filename: "config.yml",
				Applications: map[string]AppConfiguration{
					"acme-docs": {
						BlockConfiguration: BlockConfiguration{
							Type:                BlockTypeApplication,
							Name:                "acme-docs",
							ModuleSource:        "nullstone/aws-fargate-service",
							ModuleSourceVersion: latest,
							Variables: map[string]any{
								"num_tasks": 2,
							},
							Connections: map[string]types.ConnectionTarget{
								"cluster-namespace": {
									BlockName: "namespace0",
								},
							},
						},
						EnvVariables: map[string]string{
							"TESTING": "abc123",
							"BLAH":    "blahblahblah",
						},
						Capabilities: []CapabilityConfiguration{
							{
								ModuleSource:        "nullstone/aws-load-balancer",
								ModuleSourceVersion: latest,
								Variables: map[string]any{
									"health_check_path": "/status",
								},
								Connections: map[string]types.ConnectionTarget{
									"subdomain": {
										BlockName: subdomainName,
									},
								},
								Namespace: &primary,
							},
						},
					},
				},
				Datastores:        map[string]DatastoreConfiguration{},
				Subdomains:        map[string]SubdomainConfiguration{},
				Domains:           map[string]DomainConfiguration{},
				Ingresses:         map[string]IngressConfiguration{},
				ClusterNamespaces: map[string]ClusterNamespaceConfiguration{},
				Clusters:          map[string]ClusterConfiguration{},
				Networks:          map[string]NetworkConfiguration{},
				Blocks:            map[string]BlockConfiguration{},
			},
			validationErrors: nil,
		},
		{
			name:     "app module missing",
			filename: "test-fixtures/config.invalid1.yml",
			want:     nil,
			validationErrors: errors.ValidationErrors{
				{
					Context: "acme/api#config.yml (apps.acme-docs.module)\n",
					Message: "Module is required",
				},
			},
		},
		{
			name:     "invalid app module",
			filename: "test-fixtures/config.invalid2.yml",
			want:     nil,
			validationErrors: errors.ValidationErrors{
				{
					Context: "acme/api#config.yml (apps.acme-docs.module)\n",
					Message: "Module (nullstone/aws-invalid-module) does not exist",
				},
			},
		},
		{
			name:     "not an app module",
			filename: "test-fixtures/config.invalid3.yml",
			want:     nil,
			validationErrors: errors.ValidationErrors{
				{
					Context: "acme/api#config.yml (apps.acme-docs.module)\n",
					Message: "Module (nullstone/aws-s3-cdn) must be app module and match the contract (app/*/*), it is defined as capability/aws/",
				},
			},
		},
		{
			name:     "invalid app variable",
			filename: "test-fixtures/config.invalid4.yml",
			want:     nil,
			validationErrors: errors.ValidationErrors{
				errors.ValidationError{
					Context: "acme/api#config.yml (apps.acme-docs.vars.service_count)\n",
					Message: "Variable does not exist on the module (nullstone/aws-fargate-service@0.0.1)",
				},
			},
		},
		{
			name:     "capability module missing",
			filename: "test-fixtures/config.invalid5.yml",
			want:     nil,
			validationErrors: errors.ValidationErrors{
				{
					Context: "acme/api#config.yml (apps.acme-docs.capabilities[0].module)\n",
					Message: "Module is required",
				},
			},
		},
		{
			name:     "invalid capability module",
			filename: "test-fixtures/config.invalid6.yml",
			want:     nil,
			validationErrors: errors.ValidationErrors{
				{
					Context: "acme/api#config.yml (apps.acme-docs.capabilities[0].module)\n",
					Message: "Module (nullstone/aws-invalid-module) does not exist",
				},
			},
		},
		{
			name:     "not a capability module",
			filename: "test-fixtures/config.invalid7.yml",
			want:     nil,
			validationErrors: errors.ValidationErrors{
				{
					Context: "acme/api#config.yml (apps.acme-docs.capabilities[0].module)\n",
					Message: "Module (nullstone/aws-s3-site) must be capability module and match the contract (capability/aws/*), it is defined as app:static-site/aws/s3",
				},
			},
		},
		{
			name:     "capability does not match app subcategory",
			filename: "test-fixtures/config.invalid8.yml",
			want:     nil,
			validationErrors: errors.ValidationErrors{
				errors.ValidationError{
					Context: "acme/api#config.yml (apps.acme-docs.capabilities[0].module)\n",
					Message: "Module (nullstone/aws-postgres-access) does not support application category (static-site)",
				},
			},
		},
		{
			name:     "invalid capability variable",
			filename: "test-fixtures/config.invalid9.yml",
			want:     nil,
			validationErrors: errors.ValidationErrors{
				errors.ValidationError{
					Context: "acme/api#config.yml (apps.acme-docs.capabilities[0].vars.database_name)\n",
					Message: "Variable does not exist on the module (nullstone/aws-load-balancer@latest)",
				},
			},
		},
		{
			name:     "capability block doesn't exist",
			filename: "test-fixtures/config.invalid10.yml",
			want:     nil,
			validationErrors: errors.ValidationErrors{
				errors.ValidationError{
					Context: "acme/api#config.yml (apps.acme-docs.capabilities[0].connections.subdomain)\n",
					Message: "Connection is invalid, block core/ns-sub-for-blah does not exist",
				},
			},
		},
		{
			name:     "block doesn't match contract for capability connection",
			filename: "test-fixtures/config.invalid11.yml",
			want:     nil,
			validationErrors: errors.ValidationErrors{
				errors.ValidationError{
					Context: "acme/api#config.yml (apps.acme-docs.capabilities[0].connections.subdomain)\n",
					Message: "Block (postgres) does not match the required contract (subdomain/aws/route53) for the capability connection",
				},
			},
		},
		{
			name:     "blockName is required",
			filename: "test-fixtures/config.invalid12.yml",
			want:     nil,
			validationErrors: errors.ValidationErrors{
				errors.ValidationError{
					Context: "acme/api#config.yml (apps.acme-docs.capabilities[0].connections.subdomain.block_name)\n",
					Message: "Connection must have a block_name to identify which block it is connected to",
				},
			},
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

			got := ConvertConfiguration("acme/api", "config.yml", *parsed)

			if test.want != nil {
				assert.Equal(t, *test.want, got)
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
			resolver := &find.ResourceResolver{
				ApiClient:    apiHub.Client(defaults.OrgName),
				CurStackId:   defaults.StackId,
				CurEnvId:     defaults.EnvId,
				StacksById:   map[int64]*find.StackResolver{defaults.StackId: sr},
				StacksByName: map[string]*find.StackResolver{"core": sr},
			}

			err = got.Validate(resolver)
			assert.Equal(t, test.validationErrors, err)
		})
	}
}
