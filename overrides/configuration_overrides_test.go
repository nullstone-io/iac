package overrides

import (
	"github.com/nullstone-io/iac/core"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestParseConfigurationOverrides(t *testing.T) {
	namespace := "secondary"
	result := &ConfigurationOverrides{
		Version: "0.1",
		Applications: map[string]ApplicationOverrides{
			"acme-api": {
				Name: "acme-api",
				Variables: map[string]any{
					"enable_versioned_assets": false,
				},
				EnvVariables: map[string]string{
					"TESTING": "abc123",
					"BLAH":    "blahblahblah",
				},
				Capabilities: CapabilityOverrides{
					{
						ModuleSource: "nullstone/aws-s3-cdn",
						Variables:    map[string]any{"enable_www": true},
						Namespace:    &namespace,
						Connections: map[string]core.ConnectionTarget{
							"subdomain": {
								StackId:   0,
								StackName: "",
								BlockId:   0,
								BlockName: "ns-sub-for-acme-docs",
								EnvId:     nil,
								EnvName:   "",
							},
						},
					},
				},
			},
		},
	}

	tests := []struct {
		name     string
		filename string
	}{
		{
			name:     "parse overrides yaml",
			filename: "test-fixtures/previews.yml",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buf, err := ioutil.ReadFile(test.filename)
			assert.NoError(t, err)

			overrides, err := ParseConfigurationOverrides(buf)
			assert.NoError(t, err)

			assert.Equal(t, result, overrides)
		})
	}
}
