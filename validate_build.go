package iac

import (
	"slices"

	"github.com/nullstone-io/iac/config"
	"github.com/nullstone-io/iac/core"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

// `build:` of an app is checked against the builder in effect and the package mode of the app's module.
// Neither is always in the file that declares `build:`:
//   - config.yml declares the whole build config, so its builder is the one it names, else the default.
//   - An overrides file names its own builder, else it takes the builder of config.yml. When config.yml
//     declares no `build:`, the overrides are laid over what the app has outside of the files.
//   - The module is the one the overrides file names, else the one config.yml names.
//
// Validate checks with what the files name. What they do not name is not checked there;
// ValidateBuild checks it once the app is known.

// validateBuilds checks `build:` of every app in every file against what the files name
func validateBuilds(input ConfigFiles) core.ValidateErrors {
	errs := core.ValidateErrors{}
	if input.Config != nil {
		for _, name := range sortedAppNames(input.Config) {
			errs = append(errs, validateBuildIn(input.Config, nil, name, config.BuildTarget{})...)
		}
	}
	for _, key := range sortedKeys(input.Overrides) {
		overrides := input.Overrides[key]
		if overrides == nil {
			continue
		}
		for _, name := range sortedAppNames(overrides) {
			errs = append(errs, validateBuildIn(input.Config, overrides, name, config.BuildTarget{})...)
		}
	}
	return errs
}

// ValidateBuild checks the `build:` that config.yml and the overrides file of env declare for block.
// base is what the app has outside of the files: the package mode of its module, and the builder and
// docker context of the build config that an overrides file is laid over when config.yml declares none.
// Pass an empty env to check config.yml alone (the block template).
func ValidateBuild(input ConfigFiles, block types.Block, env types.Environment, base config.BuildTarget) core.ValidateErrors {
	if config.BlockType(block.Type) != config.BlockTypeApplication {
		return nil
	}
	errs := validateBuildIn(input.Config, nil, block.Name, base)
	if overrides := input.GetOverrides(env); overrides != nil {
		errs = append(errs, validateBuildIn(input.Config, overrides, block.Name, base)...)
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

// validateBuildIn checks `build:` of an app in one file: overrides when given, else config.yml.
// base.PackageMode is used when it is known, since it comes from the module the app resolved to.
// base.Builder is used only when no file names or implies a builder.
func validateBuildIn(cfg, overrides *config.EnvConfiguration, appName string, base config.BuildTarget) core.ValidateErrors {
	target := base
	file, app := cfg, findApp(cfg, appName)
	if app != nil {
		if base.PackageMode == "" && app.Module != nil {
			target.PackageMode = types.DefaultBuildConfig(app.Module).PackageMode
		}
		if app.Build != nil {
			// config.yml starts from the module default, which sets the builder and the docker context
			moduleDefault := types.DefaultBuildConfig(app.Module)
			target.Builder = moduleDefault.Builder
			target.DockerContext = moduleDefault.DockerContext
			if app.Build.Builder != nil {
				target.Builder = *app.Build.Builder
			}
			if app.Build.Nullstone != nil && app.Build.Nullstone.DockerContext != nil {
				target.DockerContext = *app.Build.Nullstone.DockerContext
			}
		}
	}
	if overrides != nil {
		file, app = overrides, findApp(overrides, appName)
		if app != nil {
			if base.PackageMode == "" && app.Module != nil {
				target.PackageMode = types.DefaultBuildConfig(app.Module).PackageMode
			}
			if app.Build != nil && app.Build.Builder != nil {
				target.Builder = *app.Build.Builder
			}
			if app.Build != nil && app.Build.Nullstone != nil && app.Build.Nullstone.DockerContext != nil {
				target.DockerContext = *app.Build.Nullstone.DockerContext
			}
		}
	}
	if app == nil || app.Build == nil {
		return nil
	}

	errs := app.Build.ValidateFit(core.NewObjectPathContextKey("apps", appName), target)
	for i := range errs {
		errs[i].IacContext = file.IacContext
	}
	return errs
}

func findApp(file *config.EnvConfiguration, name string) *config.AppConfiguration {
	if file == nil {
		return nil
	}
	return file.Applications[name]
}

func sortedAppNames(file *config.EnvConfiguration) []string {
	return sortedKeys(file.Applications)
}

func sortedKeys[T any](m map[string]T) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	return keys
}
