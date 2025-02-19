package yaml

import (
	"gopkg.in/yaml.v3"
	"strings"
)

var (
	_ yaml.Marshaler   = &ConnectionConstraint{}
	_ yaml.Unmarshaler = &ConnectionConstraint{}
)

type ConnectionConstraint struct {
	StackName string `yaml:"stack_name,omitempty"`
	EnvName   string `yaml:"env_name,omitempty"`
	BlockName string `yaml:"block_name"`
}

// connectionConstraintMap is used to parse a yaml node
// We allow a scalar node (a single string) or a mapping node (each value specified)
// This allows us to parse the mapping node
type connectionConstraintMap struct {
	StackName string `yaml:"stack_name,omitempty"`
	EnvName   string `yaml:"env_name,omitempty"`
	BlockName string `yaml:"block_name"`
}

func (c *ConnectionConstraint) UnmarshalYAML(node *yaml.Node) error {
	switch node.Kind {
	case yaml.ScalarNode:
		var full string
		if err := node.Decode(&full); err != nil {
			return err
		}
		*c = ParseConnectionConstraint(full)
	case yaml.MappingNode:
		var tmp connectionConstraintMap
		if err := node.Decode(&tmp); err != nil {
			return err
		}
		c.StackName = tmp.StackName
		c.EnvName = tmp.EnvName
		c.BlockName = tmp.BlockName
		return nil
	default:
		return nil
	}
	return nil
}

func (c ConnectionConstraint) MarshalYAML() (interface{}, error) {
	node := &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!str",
		Value: c.String(),
	}
	return node, nil
}

func (c *ConnectionConstraint) String() string {
	tokens := make([]string, 0)
	if c.StackName != "" {
		tokens = append(tokens, c.StackName)
	}
	if c.StackName != "" || c.EnvName != "" {
		// always emit an env name if stack name is added
		tokens = append(tokens, c.EnvName)
	}
	// block name should always exist
	tokens = append(tokens, c.BlockName)
	return strings.Join(tokens, ".")
}

func ParseConnectionConstraint(s string) ConnectionConstraint {
	tokens := strings.Split(s, ".")
	switch len(tokens) {
	case 1:
		return ConnectionConstraint{
			BlockName: tokens[0],
		}
	case 2:
		return ConnectionConstraint{
			EnvName:   tokens[0],
			BlockName: tokens[1],
		}
	case 3:
		return ConnectionConstraint{
			StackName: tokens[0],
			EnvName:   tokens[1],
			BlockName: tokens[2],
		}
	default:
		return ConnectionConstraint{}
	}
}
