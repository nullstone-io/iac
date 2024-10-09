package core

import (
	"fmt"
	"strconv"
	"strings"
)

type IacContext struct {
	RepoName string
	Filename string
}

func (c IacContext) Context(sub YamlPathContext) string {
	return fmt.Sprintf("%s#%s (%s)", c.RepoName, c.Filename, sub.Context())
}

func NewYamlPathContext(field string, index string) YamlPathContext {
	return YamlPathContext{Field: field, Key: index}.SubField("")
}

type YamlPathContext struct {
	Path  string
	Field string
	Index *int
	Key   string
}

func (c YamlPathContext) SubField(field string) YamlPathContext {
	path := c.Path
	if c.Field != "" {
		path = c.Context()
	}
	return YamlPathContext{
		Path:  path,
		Field: field,
	}
}

func (c YamlPathContext) SubKey(field string, key string) YamlPathContext {
	path := c.Path
	if c.Field != "" {
		path = c.Context()
	}
	return YamlPathContext{
		Path:  path,
		Field: field,
		Key:   key,
	}
}

func (c YamlPathContext) SubIndex(field string, index int) YamlPathContext {
	path := c.Path
	if c.Field != "" {
		path = c.Context()
	}
	return YamlPathContext{
		Path:  path,
		Field: field,
		Index: &index,
	}
}

func (c YamlPathContext) Context() string {
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
