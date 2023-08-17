package core

import (
	"fmt"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
	"gopkg.in/yaml.v3"
	"strings"
)

var _ yaml.Unmarshaler = &ConnectionTarget{}

type ConnectionTarget types.ConnectionTarget

type connectionTargetFull struct {
	StackId   int64  `yaml:"stack_id"`
	StackName string `yaml:"stack_name"`
	EnvId     *int64 `yaml:"env_id"`
	EnvName   string `yaml:"env_name"`
	BlockName string `yaml:"block_name"`
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

type ConnectionTargets map[string]ConnectionTarget

func (s ConnectionTargets) String() string {
	result := make([]string, 0)
	for name, c := range s {
		t := types.ConnectionTarget(c)
		result = append(result, fmt.Sprintf("%s=%s", name, t.Workspace().Id()))
	}
	return strings.Join(result, ",")
}
