package yaml

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestConnectionConstraint_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    ConnectionConstraint
	}{
		{
			name: "block_name",
			content: `target:
  block_name: "block0"`,
			want: ConnectionConstraint{
				BlockName: "block0",
			},
		},
		{
			name:    "block name as string",
			content: `target: "block0"`,
			want: ConnectionConstraint{
				BlockName: "block0",
			},
		},
		{
			name:    "stack and block name as string",
			content: `target: "stack0..block0"`,
			want: ConnectionConstraint{
				StackName: "stack0",
				BlockName: "block0",
			},
		},
		{
			name:    "env and block name as string",
			content: `target: "env0.block0"`,
			want: ConnectionConstraint{
				EnvName:   "env0",
				BlockName: "block0",
			},
		},
		{
			name:    "stack, block, and env name as string",
			content: `target: "stack0.env0.block0"`,
			want: ConnectionConstraint{
				StackName: "stack0",
				EnvName:   "env0",
				BlockName: "block0",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var got struct {
				Target ConnectionConstraint `yaml:"target"`
			}
			assert.NoError(t, yaml.Unmarshal([]byte(test.content), &got))
			assert.Equal(t, test.want, got.Target)
		})
	}
}

func TestConnectionConstraint_MarshalYAML(t *testing.T) {
	tests := []struct {
		name       string
		constraint ConnectionConstraint
		want       string
	}{
		{
			name: "block_name",
			constraint: ConnectionConstraint{
				BlockName: "block0",
			},
			want: `target: block0
`,
		},
		{
			name: "env and block name",
			constraint: ConnectionConstraint{
				EnvName:   "env0",
				BlockName: "block0",
			},
			want: `target: env0.block0
`,
		},
		{
			name: "stack and block name",
			constraint: ConnectionConstraint{
				StackName: "stack0",
				BlockName: "block0",
			},
			want: `target: stack0..block0
`,
		},
		{
			name: "stack, block, and env name",
			constraint: ConnectionConstraint{
				StackName: "stack0",
				EnvName:   "env0",
				BlockName: "block0",
			},
			want: `target: stack0.env0.block0
`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			type input struct {
				Target ConnectionConstraint `yaml:"target"`
			}
			got, err := yaml.Marshal(input{Target: test.constraint})
			assert.NoError(t, err, "unexpected error")
			assert.Equal(t, test.want, string(got))
		})
	}
}
