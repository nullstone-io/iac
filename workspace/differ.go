package workspace

import (
	"fmt"

	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

// DiffWorkspaceConfig performs a difference between two WorkspaceConfig and produces a set of []types.WorkspaceChange
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
func DiffWorkspaceConfig(cur, des types.WorkspaceConfig) IndexedChanges {
	changes := IndexedChanges{}
	diffModuleVersion(changes, cur, des)
	diffVariables(changes, cur.Variables, des.Variables)
	diffEnvVariables(changes, cur.EnvVariables, des.EnvVariables)
	diffConnections(changes, cur.Connections, des.Connections)
	diffCapabilities(changes, cur, des)
	diffExtra(changes, cur, des)
	return changes
}

func DiffCapabilityConfig(cur, des types.CapabilityConfig) IndexedChanges {
	changes := IndexedChanges{}
	diffCapModuleVersion(changes, cur, des)
	diffVariables(changes, cur.Variables, des.Variables)
	diffConnections(changes, cur.Connections, des.Connections)
	diffCapNamespace(changes, cur, des)
	return changes
}

func diffModuleVersion(changes IndexedChanges, cur, des types.WorkspaceConfig) {
	if cur.SourceVersion != "" && cur.SourceVersion != des.SourceVersion {
		changes.Add(types.WorkspaceChange{
			ChangeType: types.ChangeTypeModuleVersion,
			Identifier: types.ChangeIdentifierModuleVersion,
			Action:     types.ChangeActionUpdate,
			Current: types.ModuleConfig{
				Source:         cur.Source,
				SourceVersion:  cur.SourceVersion,
				SourceToolName: cur.SourceToolName,
				Variables:      cur.Variables,
				Connections:    cur.Connections,
				Providers:      cur.Providers,
			},
			Desired: types.ModuleConfig{
				Source:         des.Source,
				SourceVersion:  des.SourceVersion,
				SourceToolName: des.SourceToolName,
				Variables:      des.Variables,
				Connections:    des.Connections,
				Providers:      des.Providers,
			},
		})
	}
}

func diffVariables(changes IndexedChanges, cur, des types.Variables) {
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

func diffEnvVariables(changes IndexedChanges, cur, des types.EnvVariables) {
	for name, b := range des {
		if a, ok := cur[name]; !ok {
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
	for name, a := range cur {
		if _, ok := des[name]; !ok {
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

func diffConnections(changes IndexedChanges, cur, des types.Connections) {
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

func diffCapabilities(changes IndexedChanges, cur, des types.WorkspaceConfig) {
	for _, b := range des.Capabilities {
		a := cur.Capabilities.FindById(b.Id)
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
	for _, a := range cur.Capabilities {
		b := des.Capabilities.FindById(a.Id)
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

func diffCapModuleVersion(changes IndexedChanges, cur, des types.CapabilityConfig) {
	if cur.SourceVersion != "" && cur.SourceVersion != des.SourceVersion {
		changes.Add(types.WorkspaceChange{
			ChangeType: types.ChangeTypeModuleVersion,
			Identifier: types.ChangeIdentifierModuleVersion,
			Action:     types.ChangeActionUpdate,
			Current: types.ModuleConfig{
				Source:        cur.Source,
				SourceVersion: cur.SourceVersion,
				Variables:     cur.Variables,
				Connections:   cur.Connections,
			},
			Desired: types.ModuleConfig{
				Source:        des.Source,
				SourceVersion: des.SourceVersion,
				Variables:     des.Variables,
				Connections:   des.Connections,
			},
		})
	}
}

func diffCapNamespace(changes IndexedChanges, cur, des types.CapabilityConfig) {
	if cur.Namespace != des.Namespace {
		changes.Add(types.WorkspaceChange{
			ChangeType: "namespace",
			Identifier: "namespace",
			Action:     types.ChangeActionUpdate,
			Current:    cur.Namespace,
			Desired:    des.Namespace,
		})
	}
}

func diffExtra(changes IndexedChanges, current, desired types.WorkspaceConfig) {
	diffExtraSubdomain(changes, current.Extra.Subdomain, desired.Extra.Subdomain)
}

func diffExtraSubdomain(changes IndexedChanges, current, desired *types.ExtraSubdomainConfig) {
	var cur, des types.ExtraSubdomainConfig
	// Assume a nil config is a blank config
	if current != nil {
		cur = *current
	}
	if desired != nil {
		des = *desired
	}
	if cur.Equal(des) {
		return
	}

	changes.Add(types.WorkspaceChange{
		Action:     types.ChangeActionUpdate,
		ChangeType: types.ChangeTypeExtraSubdomain,
		Identifier: "extra_subdomain",
		Current:    current,
		Desired:    desired,
	})
}
