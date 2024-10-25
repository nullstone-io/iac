package workspace

import (
	"fmt"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

func DiffModuleConfig(cur, des types.ModuleConfig) IndexedChanges {
	changes := IndexedChanges{}
	diffModuleConfigVersion(changes, cur, des)
	diffVariablesSchema(changes, cur.Variables, des.Variables)
	diffConnectionsSchema(changes, cur.Connections, des.Connections)
	return changes
}

func diffModuleConfigVersion(changes IndexedChanges, cur, des types.ModuleConfig) {
	isModuleDiff := cur.Source != "" && cur.Source != des.Source
	isVersionDiff := cur.SourceVersion != "" && cur.SourceVersion != des.SourceVersion
	if !isModuleDiff && !isVersionDiff {
		return
	}
	changes.Add(types.WorkspaceChange{
		ChangeType: types.ChangeTypeModuleVersion,
		Identifier: types.ChangeIdentifierModuleVersion,
		Action:     types.ChangeActionUpdate,
		Current:    fmt.Sprintf("%s@%s", cur.Source, cur.SourceVersion),
		Desired:    fmt.Sprintf("%s@%s", des.Source, des.SourceVersion),
	})
}

func diffVariablesSchema(changes IndexedChanges, cur, des types.Variables) {
	for name, b := range des {
		// we don't need to worry about service_env_vars/env_vars or service_secrets/secrets, these are dealt with differently
		if name == "service_env_vars" || name == "service_secrets" || name == "env_vars" || name == "secrets" {
			continue
		}
		if a, ok := cur[name]; !ok {
			// add variable
			changes.Add(types.WorkspaceChange{
				ChangeType: types.ChangeTypeVariable,
				Identifier: name,
				Action:     types.ChangeActionAdd,
				Current:    nil,
				Desired:    b,
			})
		} else if !a.SchemaEquals(b) {
			// update variable
			changes.Add(types.WorkspaceChange{
				ChangeType: types.ChangeTypeVariable,
				Identifier: name,
				Action:     types.ChangeActionUpdate,
				Current:    a,
				Desired:    b,
			})
		}
	}
	for name, a := range cur {
		if _, ok := des[name]; !ok {
			// delete variable
			changes.Add(types.WorkspaceChange{
				ChangeType: types.ChangeTypeVariable,
				Identifier: name,
				Action:     types.ChangeActionDelete,
				Current:    a,
				Desired:    nil,
			})
		}
	}
}

func diffConnectionsSchema(changes IndexedChanges, cur, des types.Connections) {
	for name, b := range des {
		if a, ok := cur[name]; !ok {
			// add connection
			changes.Add(types.WorkspaceChange{
				ChangeType: types.ChangeTypeConnection,
				Identifier: name,
				Action:     types.ChangeActionAdd,
				Current:    nil,
				Desired:    b,
			})
		} else if !a.SchemaEquals(b) {
			// update connection
			changes.Add(types.WorkspaceChange{
				ChangeType: types.ChangeTypeConnection,
				Identifier: name,
				Action:     types.ChangeActionUpdate,
				Current:    a,
				Desired:    b,
			})
		}
	}
	for name, a := range cur {
		if _, ok := des[name]; !ok {
			// delete connection
			changes.Add(types.WorkspaceChange{
				ChangeType: types.ChangeTypeConnection,
				Identifier: name,
				Action:     types.ChangeActionDelete,
				Current:    a,
				Desired:    nil,
			})
		}
	}
}
