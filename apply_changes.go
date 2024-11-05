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
