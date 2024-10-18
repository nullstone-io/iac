package yaml

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestCapabilityConfigurations_UnmarshalYAML(t *testing.T) {
	tests := map[string]struct {
		content string
		want    CapabilityConfigurations
	}{
		"empty": {
			content: `~`,
			want:    nil,
		},
		"sequence": {
			content: `- module: nullstone/jwt-keys`,
			want: CapabilityConfigurations{
				{
					ModuleSource: "nullstone/jwt-keys",
				},
			},
		},
		"map": {
			content: `keys:
  module: nullstone/jwt-keys`,
			want: CapabilityConfigurations{
				{
					Name:         "keys",
					ModuleSource: "nullstone/jwt-keys",
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var got CapabilityConfigurations
			require.NoError(t, yaml.Unmarshal([]byte(test.content), &got), "unexpected error")
			require.Equal(t, test.want, got)
		})
	}
}

func TestCapabilityConfigurations_MarshalYAML(t *testing.T) {
	tests := map[string]struct {
		caps CapabilityConfigurations
		want string
	}{
		"nil": {
			caps: nil,
			want: "null\n",
		},
		"map": {
			caps: CapabilityConfigurations{
				{
					Name:         "cap1",
					ModuleSource: "nullstone/jwt-keys",
				},
				{
					Name:         "cap2",
					ModuleSource: "nullstone/rails-cookies",
				},
			},
			want: `cap1:
    module: nullstone/jwt-keys
cap2:
    module: nullstone/rails-cookies
`,
		},
		"mixed": {
			caps: CapabilityConfigurations{
				{
					Name:         "cap1",
					ModuleSource: "nullstone/jwt-keys",
				},
				{
					ModuleSource: "nullstone/rails-cookies",
				},
			},
			want: `- name: cap1
  module: nullstone/jwt-keys
- module: nullstone/rails-cookies
`,
		},
		"sequence": {
			caps: CapabilityConfigurations{
				{
					ModuleSource: "nullstone/jwt-keys",
				},
				{
					ModuleSource: "nullstone/rails-cookies",
				},
			},
			want: `- module: nullstone/jwt-keys
- module: nullstone/rails-cookies
`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := yaml.Marshal(test.caps)
			require.NoError(t, err, "unexpected error")
			assert.Equal(t, test.want, string(got))
		})
	}
}
