package events

import (
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

type ChangeAction string

const (
	ChangeActionAdd    ChangeAction = "add"
	ChangeActionUpdate ChangeAction = "update"
	ChangeActionDelete ChangeAction = "delete"
)

type Changes map[string]Change

type Change struct {
	Action  ChangeAction    `json:"action"`
	Current *types.EnvEvent `json:"current,omitempty"`
	Desired *types.EnvEvent `json:"desired,omitempty"`
}

func Diff(current, desired map[string]types.EnvEvent, initiatingRepoUrl string) Changes {
	changes := Changes{}
	if desired == nil {
		return changes
	}

	for name, des := range desired {
		des := des
		des.OwningRepoUrl = initiatingRepoUrl
		if existing, ok := current[name]; !ok {
			changes[name] = Change{
				Action:  ChangeActionAdd,
				Current: nil,
				Desired: &des,
			}
		} else if !des.IsEqual(existing) {
			changes[name] = Change{
				Action:  ChangeActionUpdate,
				Current: &existing,
				Desired: &des,
			}
		}
	}

	for name, cur := range current {
		if _, ok := desired[name]; !ok {
			changes[name] = Change{
				Action:  ChangeActionDelete,
				Current: &cur,
				Desired: nil,
			}
		}
	}

	return changes
}
