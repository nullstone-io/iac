package overrides

import (
	"github.com/nullstone-io/iac/yaml"
	"github.com/stretchr/testify/assert"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
	"io/ioutil"
	"testing"
)

func TestParseEnvOverrides(t *testing.T) {
	namespace := "secondary"
	result := &EnvOverrides{
		Applications: map[string]AppOverrides{
			"acme-api": {
				BlockOverrides: BlockOverrides{
					Name: "acme-api",
					Variables: map[string]any{
						"enable_versioned_assets": false,
					},
				},
				EnvVariables: map[string]string{
					"TESTING": "abc123",
					"BLAH":    "blahblahblah",
				},
				Capabilities: CapabilityOverrides{
					{
						ModuleSource:        "nullstone/aws-s3-cdn",
						ModuleSourceVersion: "latest",
						Variables:           map[string]any{"enable_www": true},
						Namespace:           &namespace,
						Connections: map[string]types.ConnectionTarget{
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
		Datastores:        map[string]DatastoreOverrides{},
		Subdomains:        map[string]SubdomainOverrides{},
		Domains:           map[string]DomainOverrides{},
		Ingresses:         map[string]IngressOverrides{},
		ClusterNamespaces: map[string]ClusterNamespaceOverrides{},
		Clusters:          map[string]ClusterOverrides{},
		Networks:          map[string]NetworkOverrides{},
		Blocks:            map[string]BlockOverrides{},
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

			parsed, err := yaml.ParseEnvOverrides(buf)
			assert.NoError(t, err)

			got, err := ConvertOverrides(*parsed)
			assert.NoError(t, err)

			assert.Equal(t, result, got)
		})
	}
}
