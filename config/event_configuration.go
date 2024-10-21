package config

import (
	"github.com/nullstone-io/iac/core"
	"github.com/nullstone-io/iac/yaml"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

func convertEventConfigurations(parsed yaml.EventConfigurations) EventConfigurations {
	events := EventConfigurations{}
	for name, value := range parsed {
		events[name] = eventConfigFromYaml(name, value)
	}
	return events
}

type EventConfigurations map[string]*EventConfiguration

type EventConfiguration struct {
	Name       string                    `json:"name"`
	Actions    []types.EventAction       `json:"actions"`
	BlockNames []string                  `json:"blockNames"`
	Statuses   []types.EventStatus       `json:"statuses"`
	Targets    EventTargetConfigurations `json:"targets"`

	Blocks []int64 `json:"blocks"`
}

func (c *EventConfiguration) Validate(ic core.IacContext, pc core.ObjectPathContext) core.ValidateErrors {
	return c.Targets.Validate(ic, pc)
}

func eventConfigFromYaml(name string, value yaml.EventConfiguration) *EventConfiguration {
	return &EventConfiguration{
		Name:       name,
		Actions:    value.Actions,
		BlockNames: value.Blocks,
		Statuses:   value.Statuses,
		Targets:    convertEventTargetConfigurations(value.Targets),
	}
}
