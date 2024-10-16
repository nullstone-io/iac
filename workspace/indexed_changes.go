package workspace

import (
	"fmt"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

type IndexedChanges map[string]*types.WorkspaceChange

func (i IndexedChanges) ToSlice() []types.WorkspaceChange {
	s := make([]types.WorkspaceChange, 0)
	for _, v := range i {
		s = append(s, *v)
	}
	return s
}

func (i IndexedChanges) Index(changeType types.ChangeType, identifier string) string {
	return fmt.Sprintf("%s/%s", changeType, identifier)
}

func (i IndexedChanges) Add(change types.WorkspaceChange) {
	i[i.Index(change.ChangeType, change.Identifier)] = &change
}

func (i IndexedChanges) Find(changeType types.ChangeType, identifier string) (*types.WorkspaceChange, bool) {
	c, ok := i[i.Index(changeType, identifier)]
	return c, ok
}

func (i IndexedChanges) Merge(latest IndexedChanges) IndexedChanges {
	result := IndexedChanges{}

	// Loop through mine (keep any not in "latest", adjust if in "latest")
	for k, myChange := range i {
		if otherChange, ok := latest[k]; ok {
			if merged := myChange.Merge(otherChange); merged != nil {
				result[k] = merged
			}
		} else {
			// mine has a change that's not in "latest"
			result[k] = myChange
		}
	}

	// Loop through "latest" (adjust result if not in mine)
	for k, otherChange := range latest {
		if _, ok := i[k]; !ok {
			// "latest" has a change that's not in mine
			result[k] = otherChange
		}
	}

	return result
}
