package yaml

// BuildConfiguration is `build:` on an app: which builder builds the app, and each builder's settings
// under a key named after the builder.
//
//	build:
//	  builder: nullstone
//	  nullstone:
//	    dockerfile: Dockerfile
//	    docker_context: .
//
// Every key is optional.
// In config.yml, a key that is left out takes the default of the app's module.
// In an overrides file (<env>.yml, previews.yml), a key that is left out keeps the value from config.yml.
//
// A setting that would not be used fails validation: settings of a builder other than the one in effect,
// a nullstone setting that the package mode of the app's module does not read,
// and the nullstone builder on a module it cannot build (zip).
type BuildConfiguration struct {
	// Builder is which builder builds the app: manual, nullstone or github-actions
	Builder *string `yaml:"builder,omitempty" json:"builder,omitempty"`

	// Nullstone holds the settings of the nullstone builder
	Nullstone *NullstoneBuildConfiguration `yaml:"nullstone,omitempty" json:"nullstone,omitempty"`

	// GithubActions holds the settings of the github-actions builder
	GithubActions *GithubActionsBuildConfiguration `yaml:"github-actions,omitempty" json:"githubActions,omitempty"`
}

// NullstoneBuildConfiguration holds the settings of the nullstone builder.
// Which settings are used depends on the package mode (dockerfile, zip, site-assets).
// The package mode is not a setting: it comes from the app's module.
type NullstoneBuildConfiguration struct {
	// Dockerfile and DockerContext apply to the dockerfile package mode
	Dockerfile    *string `yaml:"dockerfile,omitempty" json:"dockerfile,omitempty"`
	DockerContext *string `yaml:"docker_context,omitempty" json:"dockerContext,omitempty"`

	// BuildDir, InstallDepsCmd and BuildCmd apply to the site-assets package mode
	BuildDir       *string `yaml:"build_dir,omitempty" json:"buildDir,omitempty"`
	InstallDepsCmd *string `yaml:"install_deps_cmd,omitempty" json:"installDepsCmd,omitempty"`
	BuildCmd       *string `yaml:"build_cmd,omitempty" json:"buildCmd,omitempty"`

	// PublishDir applies to the zip and site-assets package modes
	PublishDir *string `yaml:"publish_dir,omitempty" json:"publishDir,omitempty"`
}

// GithubActionsBuildConfiguration holds the settings of the github-actions builder
type GithubActionsBuildConfiguration struct {
	WorkflowFilename *string `yaml:"workflow_filename,omitempty" json:"workflowFilename,omitempty"`
	InitiateWorkflow *bool   `yaml:"initiate_workflow,omitempty" json:"initiateWorkflow,omitempty"`
	PerformPush      *bool   `yaml:"perform_push,omitempty" json:"performPush,omitempty"`
	// WorkflowInputs are passed to the workflow when it is started.
	// An overrides file adds to or changes the inputs of config.yml by name.
	WorkflowInputs map[string]string `yaml:"workflow_inputs,omitempty" json:"workflowInputs,omitempty"`
}
