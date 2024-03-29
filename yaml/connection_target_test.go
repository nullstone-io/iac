package yaml

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestConnectionTarget_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    ConnectionTarget
	}{
		{
			name: "block_name",
			content: `target:
  block_name: "block0"`,
			want: ConnectionTarget{
				BlockName: "block0",
			},
		},
		{
			name:    "block name as string",
			content: `target: "block0"`,
			want: ConnectionTarget{
				BlockName: "block0",
			},
		},
		{
			name:    "stack and block name as string",
			content: `target: "stack0.block0"`,
			want: ConnectionTarget{
				StackName: "stack0",
				BlockName: "block0",
			},
		},
		{
			name:    "stack, block, and env name as string",
			content: `target: "stack0.env0.block0"`,
			want: ConnectionTarget{
				StackName: "stack0",
				EnvName:   "env0",
				BlockName: "block0",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var got struct {
				Target ConnectionTarget `yaml:"target"`
			}
			assert.NoError(t, yaml.Unmarshal([]byte(test.content), &got))
			assert.Equal(t, test.want, got.Target)
		})
	}
}
