package yaml

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
	goyaml "gopkg.in/yaml.v3"
)

// dataclassification is nested under metadata in .nullstone/config.yml.
func TestBlockConfiguration_UnmarshalMetadata(t *testing.T) {
	data := []byte("module: nullstone/aws-rds-postgres\nmetadata:\n  dataclassification: customer-content\n")

	var bc BlockConfiguration
	require.NoError(t, goyaml.Unmarshal(data, &bc))
	require.NotNil(t, bc.Metadata)
	require.NotNil(t, bc.Metadata.DataClassification)
	assert.Equal(t, "customer-content", *bc.Metadata.DataClassification)
}

func TestBlockConfiguration_UnmarshalWithoutMetadata(t *testing.T) {
	var bc BlockConfiguration
	require.NoError(t, goyaml.Unmarshal([]byte("module: nullstone/aws-rds-postgres\n"), &bc))
	assert.Nil(t, bc.Metadata)
}

func TestMetadataFromWorkspaceConfig(t *testing.T) {
	t.Run("unclassified yields an empty container", func(t *testing.T) {
		got := MetadataFromWorkspaceConfig(types.WorkspaceMetadata{})
		require.NotNil(t, got)
		assert.Nil(t, got.DataClassification)
	})

	t.Run("classified yields the level slug", func(t *testing.T) {
		got := MetadataFromWorkspaceConfig(types.WorkspaceMetadata{DataClassification: types.ClassificationRestricted})
		require.NotNil(t, got)
		require.NotNil(t, got.DataClassification)
		assert.Equal(t, "restricted", *got.DataClassification)
	})
}
