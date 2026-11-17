package iac

import (
	"path/filepath"
	"testing"

	"github.com/nullstone-io/iac/workspace"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

func TestApplyBuildChangesTo(t *testing.T) {
	const repoUrl, repoName = "https://github.com/acme/api", "acme/api"
	dir := filepath.Join("test-fixtures", "build1")
	source := func(filename string) *types.BuildConfigSource {
		return &types.BuildConfigSource{Repo: repoName, RepoUrl: repoUrl, File: filepath.Join(dir, filename)}
	}

	app := func(name string) types.Block {
		return types.Block{Type: string(types.BlockTypeApplication), Name: name}
	}
	dev := types.Environment{Type: types.EnvTypePipeline, Name: "dev"}
	prod := types.Environment{Type: types.EnvTypePipeline, Name: "prod"}
	preview := types.Environment{Type: types.EnvTypePreview, Name: "pr-42"}

	moduleDefault := types.DefaultBuildConfig(nil)
	// fallback is what the app has when no file declares build: it carries values no file names
	fallback := func() types.BuildConfig {
		config := types.DefaultBuildConfig(nil)
		config.Dockerfile = "ui-edited/Dockerfile"
		config.BuildCmd = "make ui-edited"
		return config
	}
	with := func(base types.BuildConfig, fn func(config *types.BuildConfig)) types.BuildConfig {
		fn(&base)
		return base
	}

	tests := []struct {
		name       string
		block      types.Block
		env        types.Environment
		want       types.BuildConfig
		wantSource *types.BuildConfigSource
	}{
		{
			name:  "template takes config.yml over the module default",
			block: app("api"),
			env:   types.Environment{},
			want: with(moduleDefault, func(config *types.BuildConfig) {
				config.Builder = "nullstone"
				config.Dockerfile = "api/Dockerfile"
				config.DockerContext = "api"
			}),
			wantSource: source("config.yml"),
		},
		{
			name:  "env without an overrides file takes config.yml",
			block: app("api"),
			env:   dev,
			want: with(moduleDefault, func(config *types.BuildConfig) {
				config.Builder = "nullstone"
				config.Dockerfile = "api/Dockerfile"
				config.DockerContext = "api"
			}),
			wantSource: source("config.yml"),
		},
		{
			name:  "overrides file sets only the keys it names",
			block: app("api"),
			env:   prod,
			want: with(moduleDefault, func(config *types.BuildConfig) {
				config.Builder = "nullstone"
				config.Dockerfile = "api/Dockerfile.prod"
				config.DockerContext = "api"
			}),
			wantSource: source("prod.yml"),
		},
		{
			name:  "overrides file adds workflow inputs by name",
			block: app("web"),
			env:   prod,
			want: with(moduleDefault, func(config *types.BuildConfig) {
				config.Builder = "github-actions"
				config.WorkflowFilename = ".github/workflows/web.yml"
				config.InitiateWorkflow = false
				config.WorkflowInputs = map[string]string{"target": "web", "minify": "true"}
			}),
			wantSource: source("prod.yml"),
		},
		{
			name:  "overrides file alone is laid over what the app already has",
			block: app("worker"),
			env:   prod,
			want: with(fallback(), func(config *types.BuildConfig) {
				config.Builder = "manual"
			}),
			wantSource: source("prod.yml"),
		},
		{
			name:       "no file declares build",
			block:      app("worker"),
			env:        dev,
			want:       fallback(),
			wantSource: nil,
		},
		{
			name:  "preview env takes previews.yml",
			block: app("api"),
			env:   preview,
			want: with(moduleDefault, func(config *types.BuildConfig) {
				config.Builder = "manual"
				config.Dockerfile = "api/Dockerfile"
				config.DockerContext = "api"
			}),
			wantSource: source("previews.yml"),
		},
		{
			name:       "block that is not an app",
			block:      types.Block{Type: string(types.BlockTypeDatastore), Name: "api"},
			env:        prod,
			want:       fallback(),
			wantSource: nil,
		},
		{
			name:       "app that is not in the files",
			block:      app("unknown"),
			env:        prod,
			want:       fallback(),
			wantSource: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			files, err := ParseConfigDir(repoUrl, repoName, dir)
			require.NoError(t, err, "cannot parse fixture dir")
			got := fallback()
			updater := &workspace.BuildConfigUpdater{Config: &got, ModuleDefault: moduleDefault}

			err = ApplyBuildChangesTo(*files, test.block, test.env, updater)

			require.NoError(t, err)
			assert.Equal(t, test.want, got)
			assert.Equal(t, test.wantSource, updater.Source)
			// the module default is only ever copied
			assert.Equal(t, types.DefaultBuildConfig(nil), moduleDefault)
		})
	}
}
