package config

import (
	"context"
	"github.com/nullstone-io/iac/core"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

// NormalizeConnectionTargets loops through all connection targets and does the following:
// 1. Fills all fields (Id+Name for Stack/Block/Env)
// 2. If block.IsShared, resolves the Env to the previews-shared env
func NormalizeConnectionTargets(ctx context.Context, connectionTargets types.ConnectionTargets, resolver core.ConnectionResolver) error {
	for key, connection := range connectionTargets {
		ct, err := resolver.ResolveConnection(ctx, connection)
		if err != nil {
			return err
		}
		connectionTargets[key] = ct
	}
	return nil
}
