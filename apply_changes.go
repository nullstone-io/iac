package iac

import (
	"github.com/nullstone-io/iac/core"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

func ApplyChangesTo(input ConfigFiles, block types.Block, env types.Environment, updater core.WorkspaceConfigUpdater) error {
	if input.Config != nil {
		if err := input.Config.ApplyChangesTo(block, updater); err != nil {
			return err
		}
	}
	if overrides := input.GetOverrides(env); overrides != nil {
		if err := overrides.ApplyChangesTo(block, updater); err != nil {
			return err
		}
	}
	return nil
}

// ApplyBuildChangesTo applies the `build:` that config.yml and then the overrides file of env declare for block.
// It is separate from ApplyChangesTo because a build config is not part of the workspace config.
// Pass an empty env to apply config.yml alone (the block template).
func ApplyBuildChangesTo(input ConfigFiles, block types.Block, env types.Environment, updater core.BuildConfigUpdater) error {
	if input.Config != nil {
		if err := input.Config.ApplyBuildChangesTo(block, updater); err != nil {
			return err
		}
	}
	if overrides := input.GetOverrides(env); overrides != nil {
		if err := overrides.ApplyBuildChangesTo(block, updater); err != nil {
			return err
		}
	}
	return nil
}
