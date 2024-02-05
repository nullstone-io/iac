package core

import (
	"gopkg.in/nullstone-io/go-api-client.v0/find"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

// NormalizeConnectionTargets loops through all connection targets and does the following:
// 1. Fills all fields (Id+Name for Stack/Block/Env)
// 2. If block.IsShared, resolves the Env to the previews-shared env
func NormalizeConnectionTargets(connectionTargets types.ConnectionTargets, resolver *find.ResourceResolver) error {
	for key, connection := range connectionTargets {
		ct, err := resolver.Resolve(connection)
		if err != nil {
			return err
		}
		connectionTargets[key] = ct
	}
	return nil
}
