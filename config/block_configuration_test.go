package config

import (
	"testing"

	"github.com/nullstone-io/iac/core"
	"github.com/nullstone-io/iac/workspace"
	"github.com/nullstone-io/iac/yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

// blockConfigFromYaml should mirror the YAML's nesting: metadata.dataclassification
// becomes Metadata.DataClassification (empty when absent).
func TestBlockConfiguration_metadataFromYaml(t *testing.T) {
	tests := []struct {
		name string
		meta *yaml.MetadataConfiguration
		want MetadataConfiguration
	}{
		{name: "no metadata", meta: nil, want: MetadataConfiguration{}},
		{name: "empty metadata", meta: &yaml.MetadataConfiguration{}, want: MetadataConfiguration{}},
		{
			name: "level set",
			meta: &yaml.MetadataConfiguration{DataClassification: ptr("customer-content")},
			want: MetadataConfiguration{DataClassification: types.ClassificationCustomerContent},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bc := blockConfigFromYaml("db", yaml.BlockConfiguration{
				ModuleSource: "nullstone/aws-rds-postgres",
				Metadata:     tt.meta,
			}, BlockTypeDatastore, types.CategoryDatastore)
			assert.Equal(t, tt.want, bc.Metadata)
		})
	}
}

// validateDataClassification should accept known taxonomy slugs (and empty), and
// reject unknown values with an error pointed at metadata.dataclassification.
func TestBlockConfiguration_validateDataClassification(t *testing.T) {
	pc := core.ObjectPathContext{Path: "datastores.customer-db"}
	tests := []struct {
		name  string
		level types.ClassificationLevel
		want  core.ValidateErrors
	}{
		{name: "empty is allowed", level: "", want: core.ValidateErrors(nil)},
		{name: "valid level", level: types.ClassificationRestricted, want: core.ValidateErrors(nil)},
		{
			name:  "invalid level",
			level: types.ClassificationLevel("top-secret"),
			want: core.ValidateErrors{
				{
					ObjectPathContext: core.ObjectPathContext{Path: "datastores.customer-db.metadata", Field: "dataclassification"},
					ErrorMessage:      "Invalid data classification value (top-secret), must be one of: public, operational, customer-content, restricted, critical",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bc := &BlockConfiguration{Metadata: MetadataConfiguration{DataClassification: tt.level}}
			assert.Equal(t, tt.want, bc.validateDataClassification(pc))
		})
	}
}

// ApplyChangesTo should thread the level onto the resolved WorkspaceConfig, and
// leave the existing value untouched when the IaC omits dataclassification.
func TestBlockConfiguration_ApplyChangesTo_dataClassification(t *testing.T) {
	t.Run("sets the level on the workspace config", func(t *testing.T) {
		wc := &types.WorkspaceConfig{}
		updater := workspace.ConfigUpdater{Config: wc}
		bc := &BlockConfiguration{Metadata: MetadataConfiguration{DataClassification: types.ClassificationCustomerContent}}
		require.NoError(t, bc.ApplyChangesTo(core.IacContext{}, updater))
		assert.Equal(t, types.ClassificationCustomerContent, wc.Metadata.DataClassification)
	})

	t.Run("empty level clears the existing classification", func(t *testing.T) {
		wc := &types.WorkspaceConfig{}
		wc.Metadata.DataClassification = types.ClassificationRestricted
		updater := workspace.ConfigUpdater{Config: wc}
		bc := &BlockConfiguration{Metadata: MetadataConfiguration{}}
		require.NoError(t, bc.ApplyChangesTo(core.IacContext{}, updater))
		assert.Equal(t, types.ClassificationLevel(""), wc.Metadata.DataClassification)
	})
}
