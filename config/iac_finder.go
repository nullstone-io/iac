package config

import (
	"context"

	"github.com/nullstone-io/iac/core"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

var (
	_ core.IacFinder = IacFinder{}
)

func NewIacFinder(config, overrides *EnvConfiguration, stackId int64, envId int64) *IacFinder {
	return &IacFinder{
		Config:    config,
		Overrides: overrides,
		StackId:   stackId,
		EnvId:     envId,
	}
}

type IacFinder struct {
	Config    *EnvConfiguration
	Overrides *EnvConfiguration
	StackId   int64
	EnvId     int64
}

// FindBlockModuleInIac looks for a BlockConfiguration in the iac configuration files for this IacFinder
// This expects the input connection target to be fully resolved (all fields populated)
func (r IacFinder) FindBlockModuleInIac(ctx context.Context, ct types.ConnectionTarget) *types.Module {
	if ct.BlockName == "" {
		return nil
	}

	// We cannot resolve module configs for blocks that are in a different stack
	if ct.StackId != r.StackId {
		return nil
	}
	// The target may resolve to a different environment than this session
	// (e.g. shared blocks resolve to `previews-shared` from a preview env).
	// A block's module identity is env-agnostic, so a same-stack match is still authoritative.

	if r.Config != nil {
		base := r.Config.FindBlockConfigurationByName(ct.BlockName)
		if base != nil && base.Module != nil {
			return base.Module
		}
	}
	if r.Overrides != nil {
		overrides := r.Overrides.FindBlockConfigurationByName(ct.BlockName)
		if overrides != nil && overrides.Module != nil {
			return overrides.Module
		}
	}
	return nil
}
