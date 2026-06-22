package yaml

import "gopkg.in/nullstone-io/go-api-client.v0/types"

const RedactedValue = "••••••••••"

type BlockConfiguration struct {
	ModuleSource     string                `yaml:"module" json:"module"`
	ModuleConstraint *string               `yaml:"module_version,omitempty" json:"moduleVersion"`
	Variables        map[string]any        `yaml:"vars,omitempty" json:"vars"`
	Connections      ConnectionConstraints `yaml:"connections,omitempty" json:"connections"`
	IsShared         bool                  `yaml:"is_shared,omitempty" json:"isShared"`
	// Metadata holds governance/descriptive metadata (e.g. data classification).
	Metadata *MetadataConfiguration `yaml:"metadata,omitempty" json:"metadata,omitempty"`
}

// MetadataConfiguration is the IaC representation of a block's metadata container.
type MetadataConfiguration struct {
	// DataClassification is the data sensitivity level slug (e.g. "customer-content").
	// Omitted/empty leaves the workspace unclassified.
	DataClassification *string `yaml:"dataclassification,omitempty" json:"dataclassification,omitempty"`
}

func BlockConfigurationFromWorkspaceConfig(stackId, envId int64, config types.WorkspaceConfig) BlockConfiguration {
	return BlockConfiguration{
		ModuleSource:     config.Source,
		ModuleConstraint: &config.SourceConstraint,
		Variables:        VariablesFromWorkspaceConfig(config.Variables),
		Connections:      ConnectionsFromWorkspaceConfig(stackId, envId, config.Connections),
		Metadata:         MetadataFromWorkspaceConfig(config.Metadata),
	}
}

func MetadataFromWorkspaceConfig(metadata types.WorkspaceMetadata) *MetadataConfiguration {
	mc := &MetadataConfiguration{}
	if metadata.DataClassification != "" {
		s := string(metadata.DataClassification)
		mc.DataClassification = &s
	}
	return mc
}

func VariablesFromWorkspaceConfig(variables types.Variables) map[string]any {
	result := map[string]any{}
	for name, value := range variables {
		if value.HasValue() {
			if value.Sensitive {
				result[name] = RedactedValue
			} else {
				result[name] = value.Value
			}
		}
	}
	return result
}

func ConnectionsFromWorkspaceConfig(stackId, envId int64, connections types.Connections) ConnectionConstraints {
	result := ConnectionConstraints{}
	for name, conn := range connections {
		if conn.DesiredTarget == nil {
			continue
		}
		cc := ConnectionConstraint{BlockName: conn.DesiredTarget.BlockName}
		if conn.DesiredTarget.StackId != stackId {
			cc.StackName = conn.DesiredTarget.StackName
		}
		if conn.DesiredTarget.EnvId != nil && *conn.DesiredTarget.EnvId == envId {
			cc.EnvName = conn.DesiredTarget.EnvName
		}
		result[name] = cc
	}
	return result
}
