package config

import (
	"github.com/nullstone-io/iac/core"
	"github.com/stretchr/testify/assert"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
	"testing"
)

func TestAppConfiguration_ApplyChanges(t *testing.T) {
	stack1 := types.Stack{
		IdModel: types.IdModel{Id: 1},
		Name:    "core",
	}
	env1 := types.Environment{
		IdModel: types.IdModel{Id: 11},
		Name:    "dev",
	}
	latest := "latest"
	subdomainName := "ns-sub-for-acme-docs"
	subdomainBlock := types.Block{
		IdModel: types.IdModel{Id: 10},
		StackId: stack1.Id,
		Name:    "ns-sub-for-acme-docs",
	}
	desiredConfig := models.DesiredConfig{
		WorkspaceConfig: models.WorkspaceConfig{
			Source:        "aws-s3-site",
			SourceVersion: "latest",
			Variables: models.Variables{
				"enable_www": {
					Value: true,
				},
			},
			EnvVariables: models.EnvVariables{},
			Capabilities: models.CapabilityConfigs{
				{
					Source:        "aws-s3-cdn",
					SourceVersion: "latest",
					Variables: models.Variables{
						"notfound_behavior": {
							Value: map[string]any{
								"document": "404.html",
								"enabled":  true,
								"spa_mode": true,
							},
						},
					},
					Connections: models.Connections{
						"subdomain": {
							Reference: &types.ConnectionTarget{
								StackId: stack1.Id,
								BlockId: subdomainBlock.Id,
								EnvId:   &env1.Id,
							},
						},
					},
					Namespace: "",
				},
			},
		},
	}

	tests := []struct {
		name     string
		changes  AppConfiguration
		expected *models.DesiredConfig
	}{
		{
			name: "applies variable, env variable, and capability variable updates",
			changes: AppConfiguration{
				Variables: map[string]any{
					"enable_www": false,
				},
				EnvVariables: map[string]string{
					"FOO": "bar",
				},
				Capabilities: []CapabilityConfiguration{
					{
						ModuleSource:        "aws-s3-cdn",
						ModuleSourceVersion: &latest,
						Variables: map[string]any{
							"notfound_behavior": map[string]any{
								"document": "404.html",
								"enabled":  true,
								"spa_mode": false,
							},
						},
						Connections: map[string]core.ConnectionTarget{
							"subdomain": {
								BlockName: subdomainName,
								StackId:   0,
								EnvId:     nil,
							},
						},
						Namespace: nil,
					},
				},
			},
			expected: &models.DesiredConfig{
				WorkspaceConfig: models.WorkspaceConfig{
					Source:        "aws-s3-site",
					SourceVersion: "latest",
					Variables: models.Variables{
						"enable_www": {
							Value: false,
						},
					},
					EnvVariables: models.EnvVariables{
						"FOO": {
							Value: "bar",
						},
					},
					Capabilities: models.CapabilityConfigs{
						{
							Source:        "aws-s3-cdn",
							SourceVersion: "latest",
							Variables: models.Variables{
								"notfound_behavior": {
									Value: map[string]any{
										"document": "404.html",
										"enabled":  true,
										"spa_mode": false,
									},
								},
							},
							Connections: models.Connections{
								"subdomain": {
									Reference: &types.ConnectionTarget{
										StackId: stack1.Id,
										BlockId: subdomainBlock.Id,
										EnvId:   &env1.Id,
									},
								},
							},
							Namespace: "",
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sr := &find.StackResolver{
				Stack:               stack1,
				PreviewsSharedEnvId: 0,
				EnvsByName:          map[string]types.Environment{env1.Name: env1},
				EnvsById:            map[int64]types.Environment{env1.Id: env1},
				BlocksById:          map[int64]types.Block{subdomainBlock.Id: subdomainBlock},
				BlocksByName:        map[string]types.Block{subdomainBlock.Name: subdomainBlock},
			}
			resolver := &find.ResourceResolver{
				ApiClient:       nil,
				CurStackId:      stack1.Id,
				CurEnvId:        env1.Id,
				CurProviderType: "aws",
				StacksById:      map[int64]*find.StackResolver{stack1.Id: sr},
				StacksByName:    map[string]*find.StackResolver{stack1.Name: sr},
			}
			assert.NoError(t, test.changes.Normalize(resolver), "unexpected error")
			err := test.changes.ApplyChangesTo(&desiredConfig)
			assert.NoError(t, err)

			assert.Equal(t, *test.expected, desiredConfig)
		})
	}
}
