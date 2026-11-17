package workspace

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

func TestDiffBuildConfig(t *testing.T) {
	base := types.DefaultBuildConfig(nil)

	tests := []struct {
		name string
		cur  types.BuildConfig
		des  types.BuildConfig
		want []types.BuildConfigChange
	}{
		{
			name: "no difference",
			cur:  base,
			des:  types.DefaultBuildConfig(nil),
			want: []types.BuildConfigChange{},
		},
		{
			name: "nil and empty workflow inputs are the same",
			cur:  types.BuildConfig{WorkflowInputs: nil},
			des:  types.BuildConfig{WorkflowInputs: map[string]string{}},
			want: []types.BuildConfigChange{},
		},
		{
			name: "build-time env vars are not reported",
			cur:  types.BuildConfig{BuildEnvVars: map[string]string{"A": "1"}},
			des:  types.BuildConfig{},
			want: []types.BuildConfigChange{},
		},
		{
			name: "every setting, in declaration order",
			cur: types.BuildConfig{
				Builder:          "nullstone",
				PackageMode:      "dockerfile",
				Dockerfile:       "Dockerfile",
				DockerContext:    ".",
				BuildDir:         "./",
				InstallDepsCmd:   "npm ci",
				BuildCmd:         "npm run build",
				PublishDir:       "./dist",
				WorkflowFilename: ".github/workflows/build.yml",
				InitiateWorkflow: true,
				PerformPush:      false,
			},
			des: types.BuildConfig{
				Builder:          "github-actions",
				PackageMode:      "site-assets",
				Dockerfile:       "web/Dockerfile",
				DockerContext:    "web",
				BuildDir:         "./web",
				InstallDepsCmd:   "yarn",
				BuildCmd:         "yarn build",
				PublishDir:       "./web/dist",
				WorkflowFilename: ".github/workflows/web.yml",
				InitiateWorkflow: false,
				PerformPush:      true,
				WorkflowInputs:   map[string]string{"target": "web"},
			},
			want: []types.BuildConfigChange{
				{Key: "builder", Current: "nullstone", Desired: "github-actions"},
				{Key: "package_mode", Current: "dockerfile", Desired: "site-assets"},
				{Key: "nullstone.dockerfile", Current: "Dockerfile", Desired: "web/Dockerfile"},
				{Key: "nullstone.docker_context", Current: ".", Desired: "web"},
				{Key: "nullstone.build_dir", Current: "./", Desired: "./web"},
				{Key: "nullstone.install_deps_cmd", Current: "npm ci", Desired: "yarn"},
				{Key: "nullstone.build_cmd", Current: "npm run build", Desired: "yarn build"},
				{Key: "nullstone.publish_dir", Current: "./dist", Desired: "./web/dist"},
				{Key: "github-actions.workflow_filename", Current: ".github/workflows/build.yml", Desired: ".github/workflows/web.yml"},
				{Key: "github-actions.initiate_workflow", Current: true, Desired: false},
				{Key: "github-actions.perform_push", Current: false, Desired: true},
				{Key: "github-actions.workflow_inputs", Current: map[string]string{}, Desired: map[string]string{"target": "web"}},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, DiffBuildConfig(test.cur, test.des))
		})
	}
}

func TestBuildConfigUpdater_ResetCopiesTheModuleDefault(t *testing.T) {
	moduleDefault := types.BuildConfig{Builder: "nullstone", WorkflowInputs: map[string]string{"a": "1"}}
	config := types.BuildConfig{Builder: "manual"}
	updater := &BuildConfigUpdater{Config: &config, ModuleDefault: moduleDefault}

	updater.ResetBuildConfig()
	config.WorkflowInputs["b"] = "2"

	assert.Equal(t, "nullstone", config.Builder)
	assert.Equal(t, map[string]string{"a": "1"}, moduleDefault.WorkflowInputs)
	assert.Nil(t, updater.Source)
}
