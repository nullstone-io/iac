package core

import (
	"gopkg.in/nullstone-io/go-api-client.v0/find"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

// NormalizeConnections loops through all connections and does the following:
// 1. Fills all fields (Id+Name for Stack/Block/Env)
// 2. If block.IsShared, resolves the Env to the previews-shared env
func NormalizeConnections(connections models.Connections, resolver *find.ResourceResolver) error {
	for key, connection := range connections {
		if connection.Reference != nil {
			reference, err := resolver.Resolve(*connection.Reference)
			if err != nil {
				return err
			}
			connection.Reference = &reference
			connections[key] = connection
		}
	}
	return nil
}

// NormalizeCapabilityConnections loops through capabilities and runs NormalizeConnections
func NormalizeCapabilityConnections(capabilities models.CapabilityConfigs, resolver *find.ResourceResolver) error {
	result := make(models.CapabilityConfigs, len(capabilities))
	for i, capability := range capabilities {
		if err := NormalizeConnections(capability.Connections, resolver); err != nil {
			return err
		}
		result[i] = capability
	}
	return nil
}

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
