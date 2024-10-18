package yaml

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

var (
	_ yaml.Unmarshaler = &CapabilityConfigurations{}
	_ yaml.Marshaler   = &CapabilityConfigurations{}
)

type CapabilityConfigurations []CapabilityConfiguration

//goland:noinspection GoMixedReceiverTypes
func (c CapabilityConfigurations) MarshalYAML() (interface{}, error) {
	if c == nil {
		return nil, nil
	}

	if c.hasAnyWithNoName() {
		return []CapabilityConfiguration(c), nil
	}

	m := map[string]CapabilityConfiguration{}
	for _, cc := range c {
		name := cc.Name
		cc.Name = ""
		m[name] = cc
	}
	return m, nil
}

//goland:noinspection GoMixedReceiverTypes
func (c *CapabilityConfigurations) UnmarshalYAML(node *yaml.Node) error {
	if node == nil {
		return nil
	}
	ccs := make([]CapabilityConfiguration, 0)
	switch node.Kind {
	default:
		return fmt.Errorf("cannot unmarshal YAML %v into CapabilityConfigurations", node.Kind)
	case yaml.SequenceNode:
		if err := node.Decode(&ccs); err != nil {
			return err
		}
	case yaml.MappingNode:
		m := map[string]CapabilityConfiguration{}
		if err := node.Decode(&m); err != nil {
			return err
		}
		for k, v := range m {
			v.Name = k
			ccs = append(ccs, v)
		}
	}
	*c = ccs
	return nil
}

//goland:noinspection GoMixedReceiverTypes
func (c CapabilityConfigurations) hasAnyWithNoName() bool {
	for _, cc := range c {
		if cc.Name == "" {
			return true
		}
	}
	return false
}

type CapabilityConfiguration struct {
	Name             string                `yaml:"name,omitempty" json:"name,omitempty"`
	ModuleSource     string                `yaml:"module,omitempty" json:"module"`
	ModuleConstraint *string               `yaml:"module_version,omitempty" json:"moduleVersion"`
	Variables        map[string]any        `yaml:"vars,omitempty" json:"vars"`
	Connections      ConnectionConstraints `yaml:"connections,omitempty" json:"connections"`
	Namespace        *string               `yaml:"namespace,omitempty" json:"namespace"`
}
