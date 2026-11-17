package config

import (
	"path"
	"slices"
	"strings"

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

// validateBuildPath ensures a path that the file names is set and stays inside the repository
func validateBuildPath(pc core.ObjectPathContext, value *string) core.ValidateErrors {
	if value == nil {
		return nil
	}
	if !isPathInsideRepo(*value) {
		return core.ValidateErrors{core.InvalidBuildPathError(pc, *value)}
	}
	return nil
}

func isPathInsideRepo(value string) bool {
	if strings.TrimSpace(value) == "" {
		return false
	}
	slashed := strings.ReplaceAll(value, `\`, "/")
	if strings.HasPrefix(slashed, "/") || strings.Contains(slashed, ":") {
		// an absolute path, a drive letter or a url
		return false
	}
	cleaned := path.Clean(slashed)
	return cleaned != ".." && !strings.HasPrefix(cleaned, "../")
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
