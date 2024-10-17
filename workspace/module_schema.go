package workspace

import (
	"github.com/nullstone-io/module/config"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

type ModuleSchema config.Manifest

func (s ModuleSchema) UpdateSchema(variables types.Variables, connections types.Connections) {
	// Add Variables defined in moduleVersion, but not in DesiredConfig
	// Update Variable schema if in DesiredConfig
	for k, v := range s.Variables {
		if existing, ok := variables[k]; ok {
			existing.Variable = v
			variables[k] = existing
		} else {
			variables[k] = types.Variable{Variable: v}
		}
	}
	// Remove Variables not in input moduleVersion schema
	for k, _ := range variables {
		if _, ok := s.Variables[k]; !ok {
			delete(variables, k)
		}
	}

	// Add Connections defined in moduleVersion, but not in DesiredConfig
	// Update Connection schema if in DesiredConfig
	for k, v := range s.Connections {
		if existing, ok := connections[k]; ok {
			existing.Connection = v
			connections[k] = existing
		} else {
			connections[k] = types.Connection{Connection: v}
		}
	}
	// Remove Connections not in input moduleVersion schema
	for k, _ := range connections {
		if _, ok := s.Connections[k]; !ok {
			delete(connections, k)
		}
	}
}
