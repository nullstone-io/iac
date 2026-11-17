package workspace

import (
	"maps"

	"github.com/nullstone-io/iac/core"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

var (
	_ core.BuildConfigUpdater = &BuildConfigUpdater{}
)

// BuildConfigUpdater collects the build config that IaC files declare for one application.
//
// Config starts as what the application falls back to when no file declares `build:`
// (its block template, else the module default) and is changed in place.
// Source stays nil when no file declared `build:`; Config is then untouched.
type BuildConfigUpdater struct {
	Config *types.BuildConfig
	// ModuleDefault is the default build config of the application's module.
	// A config.yml that declares `build:` starts from it, so a key it leaves out never keeps an old value.
	ModuleDefault types.BuildConfig
	// Source is the last file that declared `build:`
	Source *types.BuildConfigSource
}

func (u *BuildConfigUpdater) ResetBuildConfig() {
	*u.Config = CloneBuildConfig(u.ModuleDefault)
}

func (u *BuildConfigUpdater) UpdateBuildConfig(ic core.IacContext, fn func(config *types.BuildConfig)) {
	fn(u.Config)
	u.Source = &types.BuildConfigSource{Repo: ic.RepoName, RepoUrl: ic.RepoUrl, File: ic.Filename}
}

// CloneBuildConfig copies a build config so that a change to the maps of the copy leaves the original alone
func CloneBuildConfig(config types.BuildConfig) types.BuildConfig {
	config.WorkflowInputs = maps.Clone(config.WorkflowInputs)
	config.BuildEnvVars = maps.Clone(config.BuildEnvVars)
	return config
}
