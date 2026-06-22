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

// blockConfigFromYaml should lift metadata.dataclassification into the config's
// DataClassification (and leave it nil when absent).
func TestBlockConfiguration_dataClassificationFromYaml(t *testing.T) {
	tests := []struct {
		name string
		meta *yaml.MetadataConfiguration
		want *types.ClassificationLevel
	}{
		{name: "no metadata", meta: nil, want: nil},
		{name: "empty metadata", meta: &yaml.MetadataConfiguration{}, want: nil},
		{
			name: "level set",
			meta: &yaml.MetadataConfiguration{DataClassification: ptr("customer-content")},
			want: ptr(types.ClassificationCustomerContent),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bc := blockConfigFromYaml("db", yaml.BlockConfiguration{
				ModuleSource: "nullstone/aws-rds-postgres",
				Metadata:     tt.meta,
			}, BlockTypeDatastore, types.CategoryDatastore)
			assert.Equal(t, tt.want, bc.DataClassification)
		})
	}
}

// validateDataClassification should accept known taxonomy slugs (and empty), and
// reject unknown values with an error pointed at metadata.dataclassification.
func TestBlockConfiguration_validateDataClassification(t *testing.T) {
	pc := core.ObjectPathContext{Path: "datastores.customer-db"}
	tests := []struct {
		name  string
		level *types.ClassificationLevel
		want  core.ValidateErrors
	}{
		{name: "absent", level: nil, want: core.ValidateErrors(nil)},
		{name: "empty is allowed", level: ptr(types.ClassificationLevel("")), want: core.ValidateErrors{}},
		{name: "valid level", level: ptr(types.ClassificationRestricted), want: core.ValidateErrors{}},
		{
			name:  "invalid level",
			level: ptr(types.ClassificationLevel("top-secret")),
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
			bc := &BlockConfiguration{DataClassification: tt.level}
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
		bc := &BlockConfiguration{DataClassification: ptr(types.ClassificationCustomerContent)}
		require.NoError(t, bc.ApplyChangesTo(core.IacContext{}, updater))
		assert.Equal(t, types.ClassificationCustomerContent, wc.Metadata.DataClassification)
	})

	t.Run("absent leaves the existing level unchanged", func(t *testing.T) {
		wc := &types.WorkspaceConfig{}
		wc.Metadata.DataClassification = types.ClassificationRestricted
		updater := workspace.ConfigUpdater{Config: wc}
		bc := &BlockConfiguration{DataClassification: nil}
		require.NoError(t, bc.ApplyChangesTo(core.IacContext{}, updater))
		assert.Equal(t, types.ClassificationRestricted, wc.Metadata.DataClassification)
	})

	t.Run("empty string clears the level", func(t *testing.T) {
		wc := &types.WorkspaceConfig{}
		wc.Metadata.DataClassification = types.ClassificationRestricted
		updater := workspace.ConfigUpdater{Config: wc}
		bc := &BlockConfiguration{DataClassification: ptr(types.ClassificationLevel(""))}
		require.NoError(t, bc.ApplyChangesTo(core.IacContext{}, updater))
		assert.Equal(t, types.ClassificationLevel(""), wc.Metadata.DataClassification)
	})
}
