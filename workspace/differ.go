package workspace

import (
	"fmt"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

type Differ struct {
	Current types.WorkspaceConfig
	Desired types.WorkspaceConfig
}

// Diff performs a difference between two WorkspaceConfig and produces a set of []types.WorkspaceChange
//
// It is the inverse function of ApplyChanges as seen here:
// ApplyChanges(current, DiffConfig(current, desired)) == current
//
// This function performs a set difference or symmetric difference of current and desired
// This is represented by AΔB or A⊖B where A=current, B=desired
// This contains objects that belong to A or B, but not their intersection
// This results in a set of workspace changes (add, remove, update)
// - add: B-A (config in desired, not in current)
// - remove: A-B (config in current, not in desired)
// - update: A⋂B (config in current and desired, value changed)
func (d Differ) Diff() IndexedChanges {
	changes := IndexedChanges{}
	d.diffModuleVersion(changes)
	d.diffVariables(changes)
	d.diffEnvVariables(changes)
	d.diffConnections(changes)
	d.diffCapabilities(changes)
	return changes
}

func (d Differ) diffModuleVersion(changes IndexedChanges) {
	if d.Current.SourceVersion != "" && d.Current.SourceVersion != d.Desired.SourceVersion {
		changes.Add(types.WorkspaceChange{
			ChangeType: types.ChangeTypeModuleVersion,
			Identifier: types.ChangeIdentifierModuleVersion,
			Action:     types.ChangeActionUpdate,
			Current: types.ModuleConfig{
				Source:        d.Current.Source,
				SourceVersion: d.Current.SourceVersion,
				Variables:     d.Current.Variables,
				Connections:   d.Current.Connections,
				Providers:     d.Current.Providers,
			},
			Desired: types.ModuleConfig{
				Source:        d.Desired.Source,
				SourceVersion: d.Desired.SourceVersion,
				Variables:     d.Desired.Variables,
				Connections:   d.Desired.Connections,
				Providers:     d.Desired.Providers,
			},
		})
	}
}

func (d Differ) diffVariables(changes IndexedChanges) {
	for name, b := range d.Desired.Variables {
		// we don't need to worry about service_env_vars/env_vars or service_secrets/secrets, these are dealt with differently
		if name == "service_env_vars" || name == "service_secrets" || name == "env_vars" || name == "secrets" {
			continue
		}
		if a, ok := d.Current.Variables[name]; !ok {
			// add variable
			changes.Add(types.WorkspaceChange{
				ChangeType: types.ChangeTypeVariable,
				Identifier: name,
				Action:     types.ChangeActionAdd,
				Current:    nil,
				Desired:    b,
			})
		} else if !a.ValueEquals(b) {
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
	for name, a := range d.Current.Variables {
		if _, ok := d.Desired.Variables[name]; !ok {
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

func (d Differ) diffEnvVariables(changes IndexedChanges) {
	for name, b := range d.Desired.EnvVariables {
		if a, ok := d.Current.EnvVariables[name]; !ok {
			// add env variable
			changes.Add(types.WorkspaceChange{
				ChangeType: types.ChangeTypeEnvVariable,
				Identifier: name,
				Action:     types.ChangeActionAdd,
				Current:    nil,
				Desired:    b,
			})
		} else if a.Value != b.Value || a.Sensitive != b.Sensitive {
			// update env variable
			changes.Add(types.WorkspaceChange{
				ChangeType: types.ChangeTypeEnvVariable,
				Identifier: name,
				Action:     types.ChangeActionUpdate,
				Current:    a,
				Desired:    b,
			})
		}
	}
	for name, a := range d.Current.EnvVariables {
		if _, ok := d.Desired.EnvVariables[name]; !ok {
			// delete env variable
			changes.Add(types.WorkspaceChange{
				ChangeType: types.ChangeTypeEnvVariable,
				Identifier: name,
				Action:     types.ChangeActionDelete,
				Current:    a,
				Desired:    nil,
			})
		}
	}
}

func (d Differ) diffConnections(changes IndexedChanges) {
	for name, b := range d.Desired.Connections {
		if a, ok := d.Current.Connections[name]; !ok {
			// add connection
			changes.Add(types.WorkspaceChange{
				ChangeType: types.ChangeTypeConnection,
				Identifier: name,
				Action:     types.ChangeActionAdd,
				Current:    nil,
				Desired:    b,
			})
		} else if !a.TargetEquals(b) {
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
	for name, a := range d.Current.Connections {
		if _, ok := d.Desired.Connections[name]; !ok {
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

func (d Differ) diffCapabilities(changes IndexedChanges) {
	for _, b := range d.Desired.Capabilities {
		a := d.Current.Capabilities.FindById(b.Id)
		if a == nil || (a.NeedsDestroyed == true && b.NeedsDestroyed == false) {
			// add capability
			changes.Add(types.WorkspaceChange{
				ChangeType: types.ChangeTypeCapability,
				Identifier: fmt.Sprintf("%d", b.Id),
				Action:     types.ChangeActionAdd,
				Current:    nil,
				Desired:    b,
			})
			continue
		}
		// if the desired state says we need to destroy a capability, add it as a delete change
		// if the d.Current state also has the capability marked as needs destroyed, this is a cleanup run and we ignore
		if b.NeedsDestroyed == true && (a != nil && a.NeedsDestroyed == false) {
			// delete capability
			changes.Add(types.WorkspaceChange{
				ChangeType: types.ChangeTypeCapability,
				Identifier: fmt.Sprintf("%d", b.Id),
				Action:     types.ChangeActionDelete,
				Current:    *a,
				Desired:    nil,
			})
			continue
		}
		if !a.Equal(b) {
			// update capability
			changes.Add(types.WorkspaceChange{
				ChangeType: types.ChangeTypeCapability,
				Identifier: fmt.Sprintf("%d", b.Id),
				Action:     types.ChangeActionUpdate,
				Current:    *a,
				Desired:    b,
			})
		}
	}
	for _, a := range d.Current.Capabilities {
		b := d.Desired.Capabilities.FindById(a.Id)
		// if a capability was added and then deleted from the desired config and never applied, delete it
		// if the existing capability was marked as needs destroyed, then we don't need to do anything here
		if b == nil && !a.NeedsDestroyed {
			changes.Add(types.WorkspaceChange{
				ChangeType: types.ChangeTypeCapability,
				Identifier: fmt.Sprintf("%d", a.Id),
				Action:     types.ChangeActionDelete,
				Current:    a,
				Desired:    nil,
			})
		}
	}
}
