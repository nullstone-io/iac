package iac

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nullstone-io/iac/workspace"
	moduleConfig "github.com/nullstone-io/module/config"
	"github.com/stretchr/testify/require"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

func TestApplyChanges(t *testing.T) {
	input1 := types.WorkspaceConfig{
		Source:        "nullstone/aws-fargate-service",
		SourceVersion: "0.1.0",
		Providers:     []string{"aws"},
		Variables: types.Variables{
			"cpu": {
				Variable: moduleConfig.Variable{Type: "number"},
				Value:    256,
			},
			"memory": {
				Variable: moduleConfig.Variable{Type: "number"},
				Value:    512,
			},
		},
		Connections: types.Connections{},
		EnvVariables: types.EnvVariables{
			"KEY1": types.EnvVariable{Value: "value1"},
		},
		Capabilities: types.CapabilityConfigs{
			types.CapabilityConfig{
				Name:          "fake-cap",
				Source:        "nullstone/fake-cap",
				SourceVersion: "0.1.2",
				Variables: types.Variables{
					"fake-var": types.Variable{
						Variable: moduleConfig.Variable{Type: "number"},
						Value:    10,
					},
				},
				Connections: types.Connections{},
			},
		},
	}

	app1 := types.Block{Type: string(types.BlockTypeApplication), Name: "app1"}
	previewEnv1 := types.Environment{Type: types.EnvTypePreview, Name: "f-123-something"}

	tests := map[string]struct {
		input           types.WorkspaceConfig
		block           types.Block
		env             types.Environment
		testFixturesDir string
		want            types.WorkspaceConfig
	}{
		"no config, no previews, should apply no changes": {
			input:           input1,
			block:           app1,
			env:             previewEnv1,
			testFixturesDir: "scenario1",
			want:            mustClone(t, input1),
		},
		"no config, has previews, should apply overrides": {
			input:           input1,
			block:           app1,
			env:             previewEnv1,
			testFixturesDir: "scenario2",
			want: types.WorkspaceConfig{
				Source:        "nullstone/aws-fargate-service",
				SourceVersion: "0.1.0",
				Providers:     []string{"aws"},
				Variables: types.Variables{
					"cpu": {
						// changed by overrides
						Variable: moduleConfig.Variable{Type: "number"},
						Value:    512,
					},
					"memory": {
						Variable: moduleConfig.Variable{Type: "number"},
						Value:    float64(512),
					},
				},
				Connections: types.Connections{},
				EnvVariables: types.EnvVariables{
					"KEY1": types.EnvVariable{Value: "value1"},
					"KEY3": types.EnvVariable{Value: "value3"},
				},
				Capabilities: types.CapabilityConfigs{
					types.CapabilityConfig{
						// changed by overrides
						Name:          "fake-cap",
						Source:        "nullstone/fake-cap",
						SourceVersion: "0.1.2",
						Variables: types.Variables{
							"fake-var": types.Variable{
								Variable: moduleConfig.Variable{Type: "number"},
								Value:    20,
							},
						},
						Connections: types.Connections{},
					},
				},
			},
		},
		"config+previews, apply config, then overrides": {
			input:           input1,
			block:           app1,
			env:             previewEnv1,
			testFixturesDir: "scenario3",
			want: types.WorkspaceConfig{
				Source:        "nullstone/aws-fargate-service",
				SourceVersion: "0.1.0",
				Providers:     []string{"aws"},
				Variables: types.Variables{
					"cpu": {
						// changed by overrides
						Variable: moduleConfig.Variable{Type: "number"},
						Value:    512,
					},
					"memory": {
						// changed by config.yml
						Variable: moduleConfig.Variable{Type: "number"},
						Value:    2048,
					},
				},
				Connections: types.Connections{},
				EnvVariables: types.EnvVariables{
					"KEY2": types.EnvVariable{Value: "value2"},
					"KEY3": types.EnvVariable{Value: "value3"},
				},
				Capabilities: types.CapabilityConfigs{
					// cleared by config.yml
				},
			},
		},
		"config, no previews, apply config": {
			input:           input1,
			block:           app1,
			env:             previewEnv1,
			testFixturesDir: "scenario4",
			want: types.WorkspaceConfig{
				Source:        "nullstone/aws-fargate-service",
				SourceVersion: "0.1.0",
				Providers:     []string{"aws"},
				Variables: types.Variables{
					"cpu": {
						Variable: moduleConfig.Variable{Type: "number"},
						Value:    1024,
					},
					"memory": {
						Variable: moduleConfig.Variable{Type: "number"},
						Value:    float64(512),
					},
				},
				Connections: types.Connections{},
				EnvVariables: types.EnvVariables{
					"KEY2": types.EnvVariable{Value: "value2"},
				},
				Capabilities: types.CapabilityConfigs{
					types.CapabilityConfig{
						// preserved in config.yml
						Name:          "fake-cap",
						Source:        "nullstone/fake-cap",
						SourceVersion: "0.1.2",
						Variables: types.Variables{
							"fake-var": types.Variable{
								Variable: moduleConfig.Variable{Type: "number"},
								Value:    10,
							},
						},
						Connections: types.Connections{},
					},
				},
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			pmr, err := ParseConfigDir("", "TestApplyChanges", filepath.Join("test-fixtures", test.testFixturesDir))
			require.NoError(t, err, "cannot test fixture dir")
			got := mustClone(t, test.input)

			updater := workspace.ConfigUpdater{
				Config: &got,
				TemplateVars: workspace.TemplateVars{
					OrgName:   test.env.OrgName,
					StackName: "apply-changes-test",
					EnvName:   test.env.Name,
					EnvIsProd: false,
				},
			}
			err = ApplyChangesTo(*pmr, test.block, test.env, updater)
			require.NoError(t, err, "unexpected error")
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("-want, +got:\n%s", diff)
			}
		})
	}
}

func mustClone[T any](t *testing.T, wc T) T {
	raw, err := json.Marshal(wc)
	require.NoError(t, err, "marshaling clone")
	var result T
	require.NoError(t, json.Unmarshal(raw, &result), "unmarshaling clone")
	return result
}
