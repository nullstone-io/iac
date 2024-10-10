package iac

import (
	"context"
	"github.com/nullstone-io/iac/config"
	"github.com/nullstone-io/iac/core"
)

func Resolve(ctx context.Context, config *config.EnvConfiguration, overrides map[string]config.EnvConfiguration, resolver core.ModuleVersionResolver) core.ResolveErrors {
	errs := core.ResolveErrors{}
	if config != nil {
		errs = append(errs, config.Resolve(ctx, resolver)...)
	}
	for key, cur := range overrides {
		errs = append(errs, cur.Resolve(ctx, resolver)...)
		overrides[key] = cur
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}
