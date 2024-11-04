package events

import (
	"github.com/nullstone-io/iac"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

func Get(input iac.ParseMapResult, env types.Environment) map[string]types.EnvEvent {
	effective := map[string]types.EnvEvent{}
	if input.Config != nil {
		input.Config.Events.MergeInto(env, effective)
	}

	name := env.Name
	if env.Type == types.EnvTypePreview {
		name = "previews"
	}
	if overrides, ok := input.Overrides[name]; ok && overrides != nil {
		overrides.Events.MergeInto(env, effective)
	}

	return effective
}