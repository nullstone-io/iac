package config

import (
	"os"
	"testing"

	"github.com/nullstone-io/iac/core"
	"github.com/nullstone-io/iac/yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

func parseBuildFixture(t *testing.T, filename string, isOverrides bool) *EnvConfiguration {
	t.Helper()
	raw, err := os.ReadFile(filename)
	require.NoError(t, err, "read fixture")
	parsed, err := yaml.ParseEnvConfiguration(raw)
	require.NoError(t, err, "parse fixture")
	return ConvertConfiguration("https://github.com/acme/api", "acme/api", filename, isOverrides, *parsed)
}

func TestBuildConfiguration_Convert(t *testing.T) {
	ptr := func(s string) *string { return &s }
	yes, no := true, false

	ec := parseBuildFixture(t, "test-fixtures/config.build.yml", false)

	tests := map[string]*BuildConfiguration{
		"api": {
			Builder: ptr("nullstone"),
			Nullstone: &yaml.NullstoneBuildConfiguration{
				Dockerfile:    ptr("api/Dockerfile"),
				DockerContext: ptr("api"),
			},
		},
		"site": {
			Builder: ptr("nullstone"),
			Nullstone: &yaml.NullstoneBuildConfiguration{
				BuildDir:       ptr("./web"),
				InstallDepsCmd: ptr("yarn install --frozen-lockfile"),
				BuildCmd:       ptr("yarn build"),
				PublishDir:     ptr("./web/dist"),
			},
		},
		"docs": {
			Builder: ptr("github-actions"),
			GithubActions: &yaml.GithubActionsBuildConfiguration{
				WorkflowFilename: ptr(".github/workflows/docs.yml"),
				InitiateWorkflow: &yes,
				PerformPush:      &no,
				WorkflowInputs:   map[string]string{"target": "docs"},
			},
		},
		"legacy": {Builder: ptr("manual")},
		// an app that does not declare build has none
		"worker": nil,
	}
	for name, want := range tests {
		t.Run(name, func(t *testing.T) {
			app, ok := ec.Applications[name]
			require.True(t, ok, "app is in the fixture")
			assert.Equal(t, want, app.Build)
			assert.Nil(t, app.Build.Validate(core.NewObjectPathContextKey("apps", name)))
		})
	}
}

func TestBuildConfiguration_Validate(t *testing.T) {
	ec := parseBuildFixture(t, "test-fixtures/config.invalid22.yml", false)
	app := ec.Applications["api"]
	require.NotNil(t, app.Build)

	got := app.Build.Validate(core.NewObjectPathContextKey("apps", "api"))

	want := core.ValidateErrors{
		{
			ObjectPathContext: core.ObjectPathContext{Path: "apps.api.build", Field: "builder"},
			ErrorMessage:      "Invalid value (jenkins), must be one of: manual, nullstone, github-actions",
		},
		{
			ObjectPathContext: core.ObjectPathContext{Path: "apps.api.build.nullstone", Field: "dockerfile"},
			ErrorMessage:      "Invalid path (/etc/Dockerfile), must be a relative path inside the repository",
		},
		{
			ObjectPathContext: core.ObjectPathContext{Path: "apps.api.build.nullstone", Field: "docker_context"},
			ErrorMessage:      "Invalid path (../outside), must be a relative path inside the repository",
		},
		{
			ObjectPathContext: core.ObjectPathContext{Path: "apps.api.build.nullstone", Field: "build_dir"},
			ErrorMessage:      "Invalid path (), must be a relative path inside the repository",
		},
		{
			ObjectPathContext: core.ObjectPathContext{Path: "apps.api.build.nullstone", Field: "publish_dir"},
			ErrorMessage:      `Invalid path (C:\dist), must be a relative path inside the repository`,
		},
		{
			ObjectPathContext: core.ObjectPathContext{Path: "apps.api.build.github-actions", Field: "workflow_filename"},
			ErrorMessage:      "Invalid path (../../build.yml), must be a relative path inside the repository",
		},
	}
	assert.Equal(t, want, got)
}

func TestIsPathInsideRepo(t *testing.T) {
	tests := map[string]bool{
		"Dockerfile":           true,
		".":                    true,
		"./":                   true,
		"./web/dist":           true,
		"api/../Dockerfile":    true,
		"":                     false,
		"   ":                  false,
		"/Dockerfile":          false,
		"..":                   false,
		"../Dockerfile":        false,
		"api/../../Dockerfile": false,
		`..\Dockerfile`:        false,
		`C:\repo\Dockerfile`:   false,
	}
	for value, want := range tests {
		t.Run(value, func(t *testing.T) {
			assert.Equal(t, want, isPathInsideRepo(value))
		})
	}
}

func TestBuildConfiguration_ApplyTo(t *testing.T) {
	ptr := func(s string) *string { return &s }
	no := false

	t.Run("sets only the keys it names", func(t *testing.T) {
		config := types.DefaultBuildConfig(nil)
		build := &BuildConfiguration{
			Nullstone: &yaml.NullstoneBuildConfiguration{Dockerfile: ptr("api/Dockerfile")},
		}

		build.ApplyTo(&config)

		want := types.DefaultBuildConfig(nil)
		want.Dockerfile = "api/Dockerfile"
		assert.Equal(t, want, config)
	})

	t.Run("sets every key", func(t *testing.T) {
		config := types.DefaultBuildConfig(nil)
		build := &BuildConfiguration{
			Builder: ptr("github-actions"),
			Nullstone: &yaml.NullstoneBuildConfiguration{
				Dockerfile:     ptr("a"),
				DockerContext:  ptr("b"),
				BuildDir:       ptr("c"),
				InstallDepsCmd: ptr("d"),
				BuildCmd:       ptr("e"),
				PublishDir:     ptr("f"),
			},
			GithubActions: &yaml.GithubActionsBuildConfiguration{
				WorkflowFilename: ptr("g"),
				InitiateWorkflow: &no,
				PerformPush:      &no,
				WorkflowInputs:   map[string]string{"h": "i"},
			},
		}

		build.ApplyTo(&config)

		assert.Equal(t, types.BuildConfig{
			Builder: "github-actions",
			// the package mode is never set from IaC; it keeps the value from the module default
			PackageMode:      "dockerfile",
			Dockerfile:       "a",
			DockerContext:    "b",
			BuildDir:         "c",
			InstallDepsCmd:   "d",
			BuildCmd:         "e",
			PublishDir:       "f",
			WorkflowFilename: "g",
			InitiateWorkflow: false,
			PerformPush:      false,
			WorkflowInputs:   map[string]string{"h": "i"},
		}, config)
	})

	t.Run("adds workflow inputs by name", func(t *testing.T) {
		config := types.BuildConfig{WorkflowInputs: map[string]string{"target": "web", "minify": "false"}}
		build := &BuildConfiguration{
			GithubActions: &yaml.GithubActionsBuildConfiguration{WorkflowInputs: map[string]string{"minify": "true"}},
		}

		build.ApplyTo(&config)

		assert.Equal(t, map[string]string{"target": "web", "minify": "true"}, config.WorkflowInputs)
	})

	t.Run("nil changes nothing", func(t *testing.T) {
		config := types.DefaultBuildConfig(nil)
		var build *BuildConfiguration

		build.ApplyTo(&config)

		assert.Equal(t, types.DefaultBuildConfig(nil), config)
	})
}
