package yaml

import (
	"fmt"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
	"gopkg.in/yaml.v3"
	"strings"
)

var (
	_ yaml.Unmarshaler = &ConnectionTarget{}
	_ yaml.Marshaler   = &ConnectionTarget{}
)

type ConnectionTarget types.ConnectionTarget

type connectionTargetFull struct {
	StackId   int64  `yaml:"stack_id,omitempty"`
	StackName string `yaml:"stack_name,omitempty"`
	EnvId     *int64 `yaml:"env_id,omitempty"`
	EnvName   string `yaml:"env_name,omitempty"`
	BlockName string `yaml:"block_name,omitempty"`
}

func (c *ConnectionTarget) UnmarshalYAML(node *yaml.Node) error {
	switch node.Kind {
	case yaml.ScalarNode:
		var full string
		if err := node.Decode(&full); err != nil {
			return err
		}
		tmp := types.ParseConnectionTarget(full)
		*c = ConnectionTarget(tmp)
	case yaml.MappingNode:
		var tmp connectionTargetFull
		if err := node.Decode(&tmp); err != nil {
			return err
		}
		c.StackName = tmp.StackName
		c.EnvName = tmp.EnvName
		c.BlockName = tmp.BlockName
		return nil
	}
	return nil
}

func (c *ConnectionTarget) MarshalYAML() (interface{}, error) {
	tokens := make([]string, 0)
	if c.StackName != "" {
		tokens = append(tokens, c.StackName)
	}
	if c.EnvName != "" {
		tokens = append(tokens, c.EnvName)
	}
	if c.BlockName != "" {
		tokens = append(tokens, c.BlockName)
	}
	node := &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!str",
		Value: strings.Join(tokens, "."),
	}
	return node, nil
}

type ConnectionTargets map[string]ConnectionTarget

func (s ConnectionTargets) String() string {
	result := make([]string, 0)
	for name, c := range s {
		t := types.ConnectionTarget(c)
		result = append(result, fmt.Sprintf("%s=%s", name, t.Workspace().Id()))
	}
	return strings.Join(result, ",")
}
