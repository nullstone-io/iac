package config

import (
	"slices"

	"github.com/nullstone-io/iac/core"
	"github.com/nullstone-io/iac/yaml"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

var (
	buildBuilders = []string{types.BuilderManual, types.BuilderNullstone, types.BuilderGithubActions}
)

// BuildConfiguration is `build:` on an app; see yaml.BuildConfiguration for the file format.
// It holds only the keys the file names, so it can be laid over another build config.
type BuildConfiguration yaml.BuildConfiguration

func convertBuildConfiguration(parsed *yaml.BuildConfiguration) *BuildConfiguration {
	return (*BuildConfiguration)(parsed)
}

// Validate checks the values the file names. It cannot check the build config as a whole
// (e.g. that a builder other than manual has a repository), since an overrides file names only part of it.
// ValidateFit checks the settings against the builder and the module of the app.
func (b *BuildConfiguration) Validate(pc core.ObjectPathContext) core.ValidateErrors {
	if b == nil {
		return nil
	}
	errs := core.ValidateErrors{}
	pc = pc.SubField("build")

	if b.Builder != nil && !slices.Contains(buildBuilders, *b.Builder) {
		errs = append(errs, core.InvalidBuildValueError(pc.SubField("builder"), *b.Builder, buildBuilders))
	}

	if ns := b.Nullstone; ns != nil {
		nspc := pc.SubField(types.BuilderNullstone)
		// dockerfile is relative to docker_context (types.BuildConfig.DockerfilePath). Checking each path on
		// its own is enough to keep the resolved path inside the repository: the context only adds depth,
		// so the resolved path can climb out only when the dockerfile climbs out by itself.
		errs = append(errs, validateBuildPath(nspc.SubField("dockerfile"), ns.Dockerfile)...)
		errs = append(errs, validateBuildPath(nspc.SubField("docker_context"), ns.DockerContext)...)
		errs = append(errs, validateBuildPath(nspc.SubField("build_dir"), ns.BuildDir)...)
		errs = append(errs, validateBuildPath(nspc.SubField("publish_dir"), ns.PublishDir)...)
	}

	if gha := b.GithubActions; gha != nil {
		ghapc := pc.SubField(types.BuilderGithubActions)
		errs = append(errs, validateBuildPath(ghapc.SubField("workflow_filename"), gha.WorkflowFilename)...)
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

// BuildTarget is what the `build:` of one file is checked against in ValidateFit.
// An empty value is not known and is not checked.
type BuildTarget struct {
	// Builder is the builder in effect for the file: the one it names, else the one it is laid over
	Builder string
	// PackageMode is the package mode of the app's module
	PackageMode string
}

var (
	// nullstoneSettingsByPackageMode lists the settings of the nullstone builder that each package mode reads
	nullstoneSettingsByPackageMode = map[string][]string{
		types.BuildPackageModeDockerfile: {"dockerfile", "docker_context"},
		types.BuildPackageModeSiteAssets: {"build_dir", "install_deps_cmd", "build_cmd", "publish_dir"},
		types.BuildPackageModeZip:        {"publish_dir"},
	}
	// packageModesNotBuiltByNullstone lists the package modes the nullstone builder cannot build
	packageModesNotBuiltByNullstone = []string{types.BuildPackageModeZip}
)

// ValidateFit checks that the settings the file names are used: they belong to the builder in effect,
// that builder can build the package mode, and the package mode reads them.
// Validate checks the values; this checks them against the app, which one file alone may not know.
func (b *BuildConfiguration) ValidateFit(pc core.ObjectPathContext, target BuildTarget) core.ValidateErrors {
	if b == nil {
		return nil
	}
	errs := core.ValidateErrors{}
	pc = pc.SubField("build")

	builder := target.Builder
	if !slices.Contains(buildBuilders, builder) {
		// Validate reports a builder that does not exist
		builder = ""
	}
	if builder != "" {
		if b.Nullstone != nil && builder != types.BuilderNullstone {
			errs = append(errs, core.BuildSectionNotUsedError(pc.SubField(types.BuilderNullstone), types.BuilderNullstone, builder))
		}
		if b.GithubActions != nil && builder != types.BuilderGithubActions {
			errs = append(errs, core.BuildSectionNotUsedError(pc.SubField(types.BuilderGithubActions), types.BuilderGithubActions, builder))
		}
		if builder == types.BuilderNullstone && slices.Contains(packageModesNotBuiltByNullstone, target.PackageMode) {
			errs = append(errs, core.BuilderCannotBuildError(pc.SubField("builder"), builder, target.PackageMode))
		}
	}

	// A nullstone section under another builder is already reported as a whole
	used, known := nullstoneSettingsByPackageMode[target.PackageMode]
	if ns := b.Nullstone; ns != nil && known && (builder == "" || builder == types.BuilderNullstone) {
		nspc := pc.SubField(types.BuilderNullstone)
		named := []struct {
			name  string
			value *string
		}{
			{"dockerfile", ns.Dockerfile},
			{"docker_context", ns.DockerContext},
			{"build_dir", ns.BuildDir},
			{"install_deps_cmd", ns.InstallDepsCmd},
			{"build_cmd", ns.BuildCmd},
			{"publish_dir", ns.PublishDir},
		}
		for _, setting := range named {
			if setting.value != nil && !slices.Contains(used, setting.name) {
				errs = append(errs, core.BuildSettingNotUsedError(nspc.SubField(setting.name), target.PackageMode))
			}
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

// validateBuildPath ensures a path that the file names is set and stays inside the repository
func validateBuildPath(pc core.ObjectPathContext, value *string) core.ValidateErrors {
	if value == nil {
		return nil
	}
	if !types.IsRepoRelativePath(*value) {
		return core.ValidateErrors{core.InvalidBuildPathError(pc, *value)}
	}
	return nil
}

// ApplyTo lays the keys the file names over config and leaves every other setting as it is.
// The package mode is never set here: it comes from the module of the app.
func (b *BuildConfiguration) ApplyTo(config *types.BuildConfig) {
	if b == nil {
		return
	}
	setString(&config.Builder, b.Builder)
	if ns := b.Nullstone; ns != nil {
		setString(&config.Dockerfile, ns.Dockerfile)
		setString(&config.DockerContext, ns.DockerContext)
		setString(&config.BuildDir, ns.BuildDir)
		setString(&config.InstallDepsCmd, ns.InstallDepsCmd)
		setString(&config.BuildCmd, ns.BuildCmd)
		setString(&config.PublishDir, ns.PublishDir)
	}
	if gha := b.GithubActions; gha != nil {
		setString(&config.WorkflowFilename, gha.WorkflowFilename)
		if gha.InitiateWorkflow != nil {
			config.InitiateWorkflow = *gha.InitiateWorkflow
		}
		if gha.PerformPush != nil {
			config.PerformPush = *gha.PerformPush
		}
		if len(gha.WorkflowInputs) > 0 {
			if config.WorkflowInputs == nil {
				config.WorkflowInputs = map[string]string{}
			}
			for name, value := range gha.WorkflowInputs {
				config.WorkflowInputs[name] = value
			}
		}
	}
}

func setString(dst *string, src *string) {
	if src != nil {
		*dst = *src
	}
}
