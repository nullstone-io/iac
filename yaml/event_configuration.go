package yaml

import "gopkg.in/nullstone-io/go-api-client.v0/types"

type EventConfigurations map[string]EventConfiguration

type EventConfiguration struct {
	Actions  []types.EventAction      `yaml:"actions" json:"actions"`
	Blocks   []string                 `yaml:"blocks" json:"blocks"`
	Statuses []types.EventStatus      `yaml:"statuses" json:"statuses"`
	Targets  EventTargetConfiguration `yaml:"targets" json:"targets"`
}

type EventTargetConfiguration struct {
	Slack *EventTargetSlackConfiguration `yaml:"slack,omitempty" json:"slack,omitempty"`
}

type EventTargetSlackConfiguration struct {
	Channels []string `yaml:"channels" json:"channels"`
}
