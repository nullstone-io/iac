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

func TestBuildConfiguration_ValidateFit(t *testing.T) {
	ptr := func(s string) *string { return &s }
	no := false
	path := func(fields ...string) core.ObjectPathContext {
		pc := core.NewObjectPathContextKey("apps", "api")
		for _, field := range fields {
			pc = pc.SubField(field)
		}
		return pc
	}
	everyNullstoneSetting := &yaml.NullstoneBuildConfiguration{
		Dockerfile:     ptr("Dockerfile"),
		DockerContext:  ptr("."),
		BuildDir:       ptr("./"),
		InstallDepsCmd: ptr("npm ci"),
		BuildCmd:       ptr("npm run build"),
		PublishDir:     ptr("./dist"),
	}
	notUsedBy := func(packageMode string, settings ...string) core.ValidateErrors {
		errs := core.ValidateErrors{}
		for _, setting := range settings {
			errs = append(errs, core.BuildSettingNotUsedError(path("build", "nullstone", setting), packageMode))
		}
		return errs
	}

	tests := []struct {
		name   string
		build  *BuildConfiguration
		target BuildTarget
		want   core.ValidateErrors
	}{
		{
			name:   "no build",
			build:  nil,
			target: BuildTarget{Builder: "nullstone", PackageMode: "dockerfile"},
			want:   nil,
		},
		{
			name:   "nothing is known about the app",
			build:  &BuildConfiguration{Nullstone: everyNullstoneSetting, GithubActions: &yaml.GithubActionsBuildConfiguration{PerformPush: &no}},
			target: BuildTarget{},
			want:   nil,
		},
		{
			name:   "a container reads the dockerfile settings",
			build:  &BuildConfiguration{Nullstone: everyNullstoneSetting},
			target: BuildTarget{Builder: "nullstone", PackageMode: "dockerfile"},
			want:   notUsedBy("dockerfile", "build_dir", "install_deps_cmd", "build_cmd", "publish_dir"),
		},
		{
			name:   "a static site reads the site settings",
			build:  &BuildConfiguration{Nullstone: everyNullstoneSetting},
			target: BuildTarget{Builder: "nullstone", PackageMode: "site-assets"},
			want:   notUsedBy("site-assets", "dockerfile", "docker_context"),
		},
		{
			name:   "a setting no builder reads for the module is reported when the builder is not known",
			build:  &BuildConfiguration{Nullstone: &yaml.NullstoneBuildConfiguration{Dockerfile: ptr("Dockerfile"), BuildCmd: ptr("make")}},
			target: BuildTarget{PackageMode: "dockerfile"},
			want:   notUsedBy("dockerfile", "build_cmd"),
		},
		{
			name:   "the nullstone builder cannot build a zip package",
			build:  &BuildConfiguration{Nullstone: &yaml.NullstoneBuildConfiguration{PublishDir: ptr("./out")}},
			target: BuildTarget{Builder: "nullstone", PackageMode: "zip"},
			want:   core.ValidateErrors{core.BuilderCannotBuildError(path("build", "builder"), "nullstone", "zip")},
		},
		{
			name:   "another builder can build a zip package",
			build:  &BuildConfiguration{Builder: ptr("github-actions"), GithubActions: &yaml.GithubActionsBuildConfiguration{PerformPush: &no}},
			target: BuildTarget{Builder: "github-actions", PackageMode: "zip"},
			want:   nil,
		},
		{
			name:   "settings of a builder that is not selected",
			build:  &BuildConfiguration{Nullstone: everyNullstoneSetting, GithubActions: &yaml.GithubActionsBuildConfiguration{PerformPush: &no}},
			target: BuildTarget{Builder: "manual", PackageMode: "dockerfile"},
			want: core.ValidateErrors{
				core.BuildSectionNotUsedError(path("build", "nullstone"), "nullstone", "manual"),
				core.BuildSectionNotUsedError(path("build", "github-actions"), "github-actions", "manual"),
			},
		},
		{
			name:   "nullstone settings under github-actions are reported once, as a section",
			build:  &BuildConfiguration{Nullstone: everyNullstoneSetting},
			target: BuildTarget{Builder: "github-actions", PackageMode: "site-assets"},
			want:   core.ValidateErrors{core.BuildSectionNotUsedError(path("build", "nullstone"), "nullstone", "github-actions")},
		},
		{
			name:   "a builder that does not exist is left to Validate",
			build:  &BuildConfiguration{Builder: ptr("jenkins"), GithubActions: &yaml.GithubActionsBuildConfiguration{PerformPush: &no}},
			target: BuildTarget{Builder: "jenkins", PackageMode: "dockerfile"},
			want:   nil,
		},
		{
			name:   "a package mode with no known settings is not checked",
			build:  &BuildConfiguration{Nullstone: everyNullstoneSetting},
			target: BuildTarget{Builder: "nullstone", PackageMode: "buildpacks"},
			want:   nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.build.ValidateFit(core.NewObjectPathContextKey("apps", "api"), test.target)

			assert.Equal(t, test.want, got)
		})
	}
}

// TestBuildConfiguration_ValidateDockerfileAgainstContext pins the rule that dockerfile is relative to
// docker_context: the resolved path (types.BuildConfig.DockerfilePath) stays inside the repository exactly
// when the dockerfile does on its own, whatever the context, so each path is checked by itself.
func TestBuildConfiguration_ValidateDockerfileAgainstContext(t *testing.T) {
	ptr := func(s string) *string { return &s }
	dockerfilePath := core.NewObjectPathContextKey("apps", "api").SubField("build").SubField("nullstone").SubField("dockerfile")

	tests := []struct {
		name          string
		dockerfile    *string
		dockerContext *string
		wantResolved  string
		want          core.ValidateErrors
	}{
		{name: "defaults", dockerfile: ptr("Dockerfile"), dockerContext: ptr("."), wantResolved: "Dockerfile", want: nil},
		{name: "dockerfile inside the context", dockerfile: ptr("docker/Dockerfile.prod"), dockerContext: ptr("web"), wantResolved: "web/docker/Dockerfile.prod", want: nil},
		{name: "dockerfile named alone", dockerfile: ptr("Dockerfile"), dockerContext: nil, wantResolved: "Dockerfile", want: nil},
		{name: "context named alone", dockerfile: nil, dockerContext: ptr("web"), wantResolved: "web/Dockerfile", want: nil},
		{
			name:          "dockerfile climbs out of the context: rejected even though web/../.. is the repo root",
			dockerfile:    ptr("../Dockerfile"),
			dockerContext: ptr("web"),
			wantResolved:  "Dockerfile",
			want:          core.ValidateErrors{core.InvalidBuildPathError(dockerfilePath, "../Dockerfile")},
		},
		{
			name:          "dockerfile climbs out of the repo through the context",
			dockerfile:    ptr("../../Dockerfile"),
			dockerContext: ptr("web"),
			wantResolved:  "../Dockerfile",
			want:          core.ValidateErrors{core.InvalidBuildPathError(dockerfilePath, "../../Dockerfile")},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			build := &BuildConfiguration{Nullstone: &yaml.NullstoneBuildConfiguration{Dockerfile: test.dockerfile, DockerContext: test.dockerContext}}

			got := build.Validate(core.NewObjectPathContextKey("apps", "api"))
			assert.Equal(t, test.want, got)

			config := types.DefaultBuildConfig(nil)
			build.ApplyTo(&config)
			assert.Equal(t, test.wantResolved, config.DockerfilePath())
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
