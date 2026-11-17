package workspace

import (
	"maps"

	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

// DiffBuildConfig lists the settings that differ between two build configs, in the order they are declared in IaC.
// Each change is keyed by the path of the setting under `build:` (e.g. builder, nullstone.dockerfile).
// package_mode is not a key of `build:`; it comes from the module of the app and is listed when it differs
// (e.g. the module changed).
// Build-time env vars are left out: they are filled in when a build starts, not configured.
func DiffBuildConfig(cur, des types.BuildConfig) []types.BuildConfigChange {
	changes := make([]types.BuildConfigChange, 0)
	add := func(key string, current, desired any) {
		changes = append(changes, types.BuildConfigChange{Key: key, Current: current, Desired: desired})
	}
	diffString := func(key, current, desired string) {
		if current != desired {
			add(key, current, desired)
		}
	}
	diffBool := func(key string, current, desired bool) {
		if current != desired {
			add(key, current, desired)
		}
	}

	diffString("builder", cur.Builder, des.Builder)
	diffString("package_mode", cur.PackageMode, des.PackageMode)
	diffString("nullstone.dockerfile", cur.Dockerfile, des.Dockerfile)
	diffString("nullstone.docker_context", cur.DockerContext, des.DockerContext)
	diffString("nullstone.build_dir", cur.BuildDir, des.BuildDir)
	diffString("nullstone.install_deps_cmd", cur.InstallDepsCmd, des.InstallDepsCmd)
	diffString("nullstone.build_cmd", cur.BuildCmd, des.BuildCmd)
	diffString("nullstone.publish_dir", cur.PublishDir, des.PublishDir)
	diffString("github-actions.workflow_filename", cur.WorkflowFilename, des.WorkflowFilename)
	diffBool("github-actions.initiate_workflow", cur.InitiateWorkflow, des.InitiateWorkflow)
	diffBool("github-actions.perform_push", cur.PerformPush, des.PerformPush)
	// a nil map and an empty map both mean "no inputs"
	if !maps.Equal(cur.WorkflowInputs, des.WorkflowInputs) {
		add("github-actions.workflow_inputs", nonNilMap(cur.WorkflowInputs), nonNilMap(des.WorkflowInputs))
	}
	return changes
}

func nonNilMap(m map[string]string) map[string]string {
	if m == nil {
		return map[string]string{}
	}
	return m
}
