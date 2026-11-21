package iac

import (
	"path/filepath"
	"testing"

	"github.com/nullstone-io/iac/config"
	"github.com/nullstone-io/iac/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

// parseBuild2 parses the build2 fixtures and gives each app the module its file names.
// Parsing alone leaves the module unresolved; Initialize resolves it in a real run.
func parseBuild2(t *testing.T) ConfigFiles {
	t.Helper()
	files, err := ParseConfigDir("https://github.com/acme/api", "acme/api", filepath.Join("test-fixtures", "build2"))
	require.NoError(t, err, "cannot parse fixture dir")

	modules := map[string]*types.Module{
		"nullstone/aws-fargate-service": {OrgName: "nullstone", Name: "aws-fargate-service", Subcategory: types.SubcategoryAppContainer},
		"nullstone/aws-s3-site":         {OrgName: "nullstone", Name: "aws-s3-site", Subcategory: types.SubcategoryAppStaticSite},
		"nullstone/aws-lambda-zip":      {OrgName: "nullstone", Name: "aws-lambda-zip", Subcategory: types.SubcategoryAppServerless, Subplatform: "zip"},
	}
	resolve := func(file *config.EnvConfiguration) {
		for _, app := range file.Applications {
			if module, ok := modules[app.ModuleSource]; ok {
				app.Module = module
				app.ModuleVersion = &types.ModuleVersion{Version: "1.0.0"}
			}
		}
	}
	resolve(files.Config)
	for _, overrides := range files.Overrides {
		resolve(overrides)
	}
	return *files
}

func buildPath(app string, fields ...string) core.ObjectPathContext {
	pc := core.NewObjectPathContextKey("apps", app)
	for _, field := range fields {
		pc = pc.SubField(field)
	}
	return pc
}

func TestValidate_Build(t *testing.T) {
	files := parseBuild2(t)
	inConfig, inProd := files.Config.IacContext, files.Overrides["prod"].IacContext

	got := Validate(files)

	want := core.ValidateErrors{
		{
			IacContext:        inConfig,
			ObjectPathContext: buildPath("fn", "build", "builder"),
			ErrorMessage:      "The nullstone builder cannot build the package mode (zip) of the app's module",
		},
		{
			IacContext:        inConfig,
			ObjectPathContext: buildPath("site", "build", "nullstone", "dockerfile"),
			ErrorMessage:      "Setting is not used by the package mode (site-assets) of the app's module",
		},
		{
			IacContext:        inConfig,
			ObjectPathContext: buildPath("web", "build", "nullstone"),
			ErrorMessage:      "Settings of the nullstone builder are not used, the builder is (github-actions)",
		},
		// prod.yml takes the module of each app from config.yml
		{
			IacContext:        inProd,
			ObjectPathContext: buildPath("site", "build", "nullstone", "docker_context"),
			ErrorMessage:      "Setting is not used by the package mode (site-assets) of the app's module",
		},
		// prod.yml takes the builder from config.yml when it names none
		{
			IacContext:        inProd,
			ObjectPathContext: buildPath("web", "build", "nullstone"),
			ErrorMessage:      "Settings of the nullstone builder are not used, the builder is (github-actions)",
		},
		// The builder of worker is not known, so its github-actions settings are not checked.
		// Its module is, so a setting no builder would read for it is reported.
		{
			IacContext:        inProd,
			ObjectPathContext: buildPath("worker", "build", "nullstone", "build_dir"),
			ErrorMessage:      "Setting is not used by the package mode (dockerfile) of the app's module",
		},
	}
	assert.Equal(t, want, got)
}

func TestValidateBuild(t *testing.T) {
	app := func(name string) types.Block {
		return types.Block{Type: string(types.BlockTypeApplication), Name: name}
	}
	prod := types.Environment{Type: types.EnvTypePipeline, Name: "prod"}
	dev := types.Environment{Type: types.EnvTypePipeline, Name: "dev"}
	files := parseBuild2(t)
	inConfig, inProd := files.Config.IacContext, files.Overrides["prod"].IacContext

	tests := []struct {
		name  string
		block types.Block
		env   types.Environment
		base  config.BuildTarget
		want  core.ValidateErrors
	}{
		{
			name:  "an app in no config.yml is checked against its own module and builder",
			block: app("cron"),
			env:   prod,
			base:  config.BuildTarget{Builder: types.BuilderNullstone, PackageMode: types.BuildPackageModeDockerfile},
			want: core.ValidateErrors{
				{
					IacContext:        inProd,
					ObjectPathContext: buildPath("cron", "build", "nullstone", "build_cmd"),
					ErrorMessage:      "Setting is not used by the package mode (dockerfile) of the app's module",
				},
			},
		},
		{
			name:  "the same settings fit a static site",
			block: app("cron"),
			env:   prod,
			base:  config.BuildTarget{Builder: types.BuilderNullstone, PackageMode: types.BuildPackageModeSiteAssets},
			want:  nil,
		},
		{
			name:  "overrides alone are checked against the builder they are laid over",
			block: app("worker"),
			env:   prod,
			base:  config.BuildTarget{Builder: types.BuilderNullstone, PackageMode: types.BuildPackageModeDockerfile},
			want: core.ValidateErrors{
				{
					IacContext:        inProd,
					ObjectPathContext: buildPath("worker", "build", "github-actions"),
					ErrorMessage:      "Settings of the github-actions builder are not used, the builder is (nullstone)",
				},
				{
					IacContext:        inProd,
					ObjectPathContext: buildPath("worker", "build", "nullstone", "build_dir"),
					ErrorMessage:      "Setting is not used by the package mode (dockerfile) of the app's module",
				},
			},
		},
		{
			name:  "the builder config.yml declares wins over the one outside of the files",
			block: app("web"),
			env:   prod,
			base:  config.BuildTarget{Builder: types.BuilderNullstone, PackageMode: types.BuildPackageModeDockerfile},
			want: core.ValidateErrors{
				{
					IacContext:        inConfig,
					ObjectPathContext: buildPath("web", "build", "nullstone"),
					ErrorMessage:      "Settings of the nullstone builder are not used, the builder is (github-actions)",
				},
				{
					IacContext:        inProd,
					ObjectPathContext: buildPath("web", "build", "nullstone"),
					ErrorMessage:      "Settings of the nullstone builder are not used, the builder is (github-actions)",
				},
			},
		},
		{
			name:  "the package mode the app resolved to wins over the module in the files",
			block: app("api"),
			env:   dev,
			base:  config.BuildTarget{PackageMode: types.BuildPackageModeSiteAssets},
			want: core.ValidateErrors{
				{
					IacContext:        inConfig,
					ObjectPathContext: buildPath("api", "build", "nullstone", "dockerfile"),
					ErrorMessage:      "Setting is not used by the package mode (site-assets) of the app's module",
				},
			},
		},
		{
			name:  "an empty env checks config.yml alone",
			block: app("site"),
			env:   types.Environment{},
			base:  config.BuildTarget{PackageMode: types.BuildPackageModeSiteAssets},
			want: core.ValidateErrors{
				{
					IacContext:        inConfig,
					ObjectPathContext: buildPath("site", "build", "nullstone", "dockerfile"),
					ErrorMessage:      "Setting is not used by the package mode (site-assets) of the app's module",
				},
			},
		},
		{
			name:  "settings that fit",
			block: app("api"),
			env:   prod,
			base:  config.BuildTarget{Builder: types.BuilderNullstone, PackageMode: types.BuildPackageModeDockerfile},
			want:  nil,
		},
		{
			name:  "block that is not an app",
			block: types.Block{Type: string(types.BlockTypeDatastore), Name: "site"},
			env:   prod,
			base:  config.BuildTarget{PackageMode: types.BuildPackageModeSiteAssets},
			want:  nil,
		},
		{
			name:  "app that is not in the files",
			block: app("unknown"),
			env:   prod,
			base:  config.BuildTarget{Builder: types.BuilderNullstone, PackageMode: types.BuildPackageModeDockerfile},
			want:  nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := ValidateBuild(files, test.block, test.env, test.base)

			assert.Equal(t, test.want, got)
		})
	}
}
