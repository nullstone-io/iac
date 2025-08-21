package workspace

import "regexp"

var (
	templateVarRegexp = map[string]*regexp.Regexp{
		"ORG":   regexp.MustCompile(`\{\{\s*NULLSTONE_ORG\s*}}`),
		"STACK": regexp.MustCompile(`\{\{\s*NULLSTONE_STACK\s*}}`),
		"ENV":   regexp.MustCompile(`\{\{\s*NULLSTONE_ENV\s*}}`),
	}
)

type TemplateVars struct {
	OrgName   string
	StackName string
	EnvName   string
}

func (v TemplateVars) ReplaceVars(input string) string {
	replaced := input
	replaced = templateVarRegexp["ORG"].ReplaceAllString(replaced, v.OrgName)
	replaced = templateVarRegexp["STACK"].ReplaceAllString(replaced, v.StackName)
	replaced = templateVarRegexp["ENV"].ReplaceAllString(replaced, v.EnvName)
	return replaced
}
