package workspace

import (
	"maps"
	"regexp"
	"slices"
)

var (
	templateVarRegexp = map[string]*regexp.Regexp{
		"NULLSTONE_ORG":   regexp.MustCompile(`\{\{\s*NULLSTONE_ORG\s*}}`),
		"NULLSTONE_STACK": regexp.MustCompile(`\{\{\s*NULLSTONE_STACK\s*}}`),
		"NULLSTONE_ENV":   regexp.MustCompile(`\{\{\s*NULLSTONE_ENV\s*}}`),
	}
)

type TemplateVars struct {
	OrgName   string
	StackName string
	EnvName   string
}

func (v TemplateVars) ReplaceVars(input string) string {
	return v.ReplaceSpecificVars(input, slices.Collect(maps.Keys(templateVarRegexp))...)
}

func (v TemplateVars) ReplaceSpecificVars(input string, vars ...string) string {
	replaced := input
	for _, r := range vars {
		re, ok := templateVarRegexp[r]
		if !ok {
			continue
		}
		replaced = re.ReplaceAllString(replaced, v.value(r))
	}
	return replaced
}

func (v TemplateVars) value(key string) string {
	switch key {
	case "NULLSTONE_ORG":
		return v.OrgName
	case "NULLSTONE_STACK":
		return v.StackName
	case "NULLSTONE_ENV":
		return v.EnvName
	default:
		return "{{ " + key + " }}"
	}
}
