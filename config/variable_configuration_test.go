package config

import (
	"testing"

	"github.com/nullstone-io/module/config"
	"gopkg.in/yaml.v3"
)

func TestExtractBaseType(t *testing.T) {
	tests := []struct {
		name     string
		varType  string
		expected string
	}{
		{
			name:     "primitive string",
			varType:  "string",
			expected: "string",
		},
		{
			name:     "primitive number",
			varType:  "number",
			expected: "number",
		},
		{
			name:     "primitive bool",
			varType:  "bool",
			expected: "bool",
		},
		{
			name:     "list of strings",
			varType:  "list(string)",
			expected: "list",
		},
		{
			name:     "map of numbers",
			varType:  "map(number)",
			expected: "map",
		},
		{
			name:     "set of strings",
			varType:  "set(string)",
			expected: "set",
		},
		{
			name:     "object with attributes",
			varType:  "object({ enabled = bool, name = string })",
			expected: "object",
		},
		{
			name:     "tuple with mixed types",
			varType:  "tuple([string, number, bool])",
			expected: "tuple",
		},
		{
			name:     "list with whitespace",
			varType:  "  list(string)  ",
			expected: "list",
		},
		{
			name:     "object with complex attributes",
			varType:  "object({ enabled : bool spa_mode : bool document : string })",
			expected: "object",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractBaseType(tt.varType)
			if got != tt.expected {
				t.Errorf("extractBaseType(%q) = %q, want %q", tt.varType, got, tt.expected)
			}
		})
	}
}

func TestIsValueCompatibleWithType(t *testing.T) {
	tests := []struct {
		name     string
		yamlVars string
		schemas  map[string]config.Variable
		expected map[string]bool // key: variable name, value: expected compatibility
	}{
		{
			name: "primitive types - valid",
			yamlVars: `
vars:
  string_var: "hello world"
  number_var: 42
  bool_var: true
`,
			schemas: map[string]config.Variable{
				"string_var": {Type: "string"},
				"number_var": {Type: "number"},
				"bool_var":   {Type: "bool"},
			},
			expected: map[string]bool{
				"string_var": true,
				"number_var": true,
				"bool_var":   true,
			},
		},
		{
			name: "primitive types - invalid",
			yamlVars: `
vars:
  string_var: 123
  number_var: "not a number"
  bool_var: "not a bool"
`,
			schemas: map[string]config.Variable{
				"string_var": {Type: "string"},
				"number_var": {Type: "number"},
				"bool_var":   {Type: "bool"},
			},
			expected: map[string]bool{
				"string_var": false,
				"number_var": false,
				"bool_var":   false,
			},
		},
		{
			name: "collection types - valid",
			yamlVars: `
vars:
  list_var: ["item1", "item2", "item3"]
  set_var: ["unique1", "unique2"]
  tuple_var: ["string", 123, true]
`,
			schemas: map[string]config.Variable{
				"list_var":  {Type: "list(string)"},
				"set_var":   {Type: "set(string)"},
				"tuple_var": {Type: "tuple([string, number, bool])"},
			},
			expected: map[string]bool{
				"list_var":  true,
				"set_var":   true,
				"tuple_var": true,
			},
		},
		{
			name: "collection types - invalid",
			yamlVars: `
vars:
  list_var: "not a list"
  set_var: 123
  tuple_var: {"not": "a tuple"}
`,
			schemas: map[string]config.Variable{
				"list_var":  {Type: "list(string)"},
				"set_var":   {Type: "set(string)"},
				"tuple_var": {Type: "tuple([string, number, bool])"},
			},
			expected: map[string]bool{
				"list_var":  false,
				"set_var":   false,
				"tuple_var": false,
			},
		},
		{
			name: "map/object types - valid",
			yamlVars: `
vars:
  map_var:
    key1: "value1"
    key2: "value2"
  object_var:
    name: "test"
    enabled: true
    count: 42
`,
			schemas: map[string]config.Variable{
				"map_var":    {Type: "map(string)"},
				"object_var": {Type: "object({ name = string, enabled = bool, count = number })"},
			},
			expected: map[string]bool{
				"map_var":    true,
				"object_var": true,
			},
		},
		{
			name: "map/object types - invalid",
			yamlVars: `
vars:
  map_var: "not a map"
  object_var: ["not", "an", "object"]
`,
			schemas: map[string]config.Variable{
				"map_var":    {Type: "map(string)"},
				"object_var": {Type: "object({ name = string, enabled = bool, count = number })"},
			},
			expected: map[string]bool{
				"map_var":    false,
				"object_var": false,
			},
		},
		{
			name: "null values - should be compatible",
			yamlVars: `
vars:
  string_var: null
  number_var: null
  bool_var: null
  list_var: null
  map_var: null
  object_var: null
`,
			schemas: map[string]config.Variable{
				"string_var": {Type: "string"},
				"number_var": {Type: "number"},
				"bool_var":   {Type: "bool"},
				"list_var":   {Type: "list(string)"},
				"map_var":    {Type: "map(string)"},
				"object_var": {Type: "object({ name = string })"},
			},
			expected: map[string]bool{
				"string_var": true,
				"number_var": true,
				"bool_var":   true,
				"list_var":   true,
				"map_var":    true,
				"object_var": true,
			},
		},
		{
			name: "numeric types - valid",
			yamlVars: `
vars:
  int_var: 42
  float_var: 3.14
  zero_var: 0
  negative_var: -123
`,
			schemas: map[string]config.Variable{
				"int_var":      {Type: "number"},
				"float_var":    {Type: "number"},
				"zero_var":     {Type: "number"},
				"negative_var": {Type: "number"},
			},
			expected: map[string]bool{
				"int_var":      true,
				"float_var":    true,
				"zero_var":     true,
				"negative_var": true,
			},
		},
		{
			name: "complex nested structures",
			yamlVars: `
vars:
  complex_list:
    - name: "item1"
      count: 10
      enabled: true
    - name: "item2"
      count: 20
      enabled: false
  complex_map:
    nested:
      deep:
        value: "found"
    numbers: [1, 2, 3]
`,
			schemas: map[string]config.Variable{
				"complex_list": {Type: "list(object({ name = string, count = number, enabled = bool }))"},
				"complex_map":  {Type: "map(any)"},
			},
			expected: map[string]bool{
				"complex_list": true, // YAML parses as []any, which is compatible with list
				"complex_map":  true, // YAML parses as map[string]any, which is compatible with map
			},
		},
		{
			name: "unknown types - should be compatible",
			yamlVars: `
vars:
  unknown_type_var: "anything"
  another_unknown: 123
`,
			schemas: map[string]config.Variable{
				"unknown_type_var": {Type: "custom_type"},
				"another_unknown":  {Type: "unknown_format"},
			},
			expected: map[string]bool{
				"unknown_type_var": true,
				"another_unknown":  true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse YAML like the actual library does
			var parsed struct {
				Vars map[string]any `yaml:"vars"`
			}
			err := yaml.Unmarshal([]byte(tt.yamlVars), &parsed)
			if err != nil {
				t.Fatalf("Failed to parse YAML: %v", err)
			}

			// Test each variable
			for varName, expectedCompatible := range tt.expected {
				t.Run(varName, func(t *testing.T) {
					value := parsed.Vars[varName]
					schema := tt.schemas[varName]
					got := isValueCompatibleWithType(value, schema.Type)
					if got != expectedCompatible {
						t.Errorf("isValueCompatibleWithType(%v (type: %T), %q) = %v, want %v",
							value, value, schema.Type, got, expectedCompatible)
					}
				})
			}
		})
	}
}
