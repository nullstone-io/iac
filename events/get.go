package events

import (
	"github.com/nullstone-io/iac"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

func Get(input iac.ConfigFiles, env types.Environment) map[string]types.EnvEvent {
	overrides := input.GetOverrides(env)
	if input.Config == nil && overrides == nil {
		return nil
	}
	effective := map[string]types.EnvEvent{}
	if input.Config != nil {
		input.Config.Events.MergeInto(env, effective)
	}
	if overrides != nil {
		overrides.Events.MergeInto(env, effective)
	}
	return effective
}
