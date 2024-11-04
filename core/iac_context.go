package core

import (
	"fmt"
	"strconv"
	"strings"
)

type IacContext struct {
	RepoUrl     string `json:"repoUrl"`
	RepoName    string `json:"repoName"`
	Filename    string `json:"filename"`
	IsOverrides bool   `json:"isOverrides"`
	Version     string `json:"version"`
}

func (c IacContext) Context(sub ObjectPathContext) string {
	return fmt.Sprintf("%s#%s (%s)", c.RepoName, c.Filename, sub.Context())
}

func NewObjectPathContextKey(field string, key string) ObjectPathContext {
	return ObjectPathContext{Field: field, Key: key}.SubField("")
}

func NewObjectPathContextIndex(field string, index int) ObjectPathContext {
	return ObjectPathContext{Field: field, Index: &index}.SubField("")
}

type ObjectPathContext struct {
	Path  string `json:"path"`
	Field string `json:"field"`
	Index *int   `json:"index,omitempty"`
	Key   string `json:"key,omitempty"`
}

func (c ObjectPathContext) SubField(field string) ObjectPathContext {
	path := c.Path
	if c.Field != "" {
		path = c.Context()
	}
	return ObjectPathContext{
		Path:  path,
		Field: field,
	}
}

func (c ObjectPathContext) SubKey(field string, key string) ObjectPathContext {
	path := c.Path
	if c.Field != "" {
		path = c.Context()
	}
	return ObjectPathContext{
		Path:  path,
		Field: field,
		Key:   key,
	}
}

func (c ObjectPathContext) SubIndex(field string, index int) ObjectPathContext {
	path := c.Path
	if c.Field != "" {
		path = c.Context()
	}
	return ObjectPathContext{
		Path:  path,
		Field: field,
		Index: &index,
	}
}

func (c ObjectPathContext) Context() string {
	started := false
	sb := strings.Builder{}
	if c.Path != "" {
		sb.WriteString(c.Path)
		started = true
	}
	if c.Field != "" {
		if started {
			sb.WriteString(".")
		}
		started = true
		sb.WriteString(c.Field)
	}
	if c.Key != "" {
		if started {
			sb.WriteString(".")
		}
		started = true
		sb.WriteString(c.Key)
	} else if c.Index != nil {
		sb.WriteString("[")
		sb.WriteString(strconv.Itoa(*c.Index))
		sb.WriteString("]")
	}
	return sb.String()
}
