package config

import (
	"context"
	"fmt"
	"slices"

	"github.com/nullstone-io/iac/core"
	"github.com/nullstone-io/iac/yaml"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

func convertEventConfigurations(parsed yaml.EventConfigurations) EventConfigurations {
	if parsed == nil {
		return nil
	}
	events := EventConfigurations{}
	for name, value := range parsed {
		events[name] = eventConfigFromYaml(name, value)
	}
	return events
}

type EventConfigurations map[string]*EventConfiguration

func (s EventConfigurations) MergeInto(env types.Environment, events map[string]types.EnvEvent) {
	for name, cur := range s {
		if existing, ok := events[name]; !ok {
			// event doesn't exist yet, add a new one
			events[name] = cur.ToEnvEvent(name, env)
		} else {
			// event exists, perform override
			events[name] = cur.OverrideEnvEvent(existing)
		}
	}
}

type EventConfiguration struct {
	Name       string                    `json:"name"`
	Actions    []types.EventAction       `json:"actions"`
	BlockNames []string                  `json:"blockNames"`
	Statuses   []types.EventStatus       `json:"statuses"`
	Targets    EventTargetConfigurations `json:"targets"`

	Blocks types.Blocks `json:"blocks"`
}

func (c *EventConfiguration) Resolve(ctx context.Context, resolver core.ResolveResolver, finder core.IacFinder, ic core.IacContext, pc core.ObjectPathContext) core.ResolveErrors {
	errs := core.ResolveErrors{}
	errs = append(errs, c.Targets.Resolve(ctx, resolver, ic, pc)...)
	return errs
}

func (c *EventConfiguration) Validate(ic core.IacContext, pc core.ObjectPathContext) core.ValidateErrors {
	errs := core.ValidateErrors{}
	if err := c.validateActions(pc); err != nil {
		errs = append(errs, *err)
	}
	if err := c.validateStatuses(pc); err != nil {
		errs = append(errs, *err)
	}
	return append(errs, c.Targets.Validate(ic, pc)...)
}

func (c *EventConfiguration) validateActions(pc core.ObjectPathContext) *core.ValidateError {
	invalid := make([]string, 0)
	for _, action := range c.Actions {
		if !slices.Contains(types.AllEventActions, action) {
			invalid = append(invalid, string(action))
		}
	}
	if len(invalid) < 1 {
		return nil
	}
	return core.InvalidEventActionError(pc, invalid)
}

func (c *EventConfiguration) validateStatuses(pc core.ObjectPathContext) *core.ValidateError {
	invalid := make([]string, 0)
	for _, status := range c.Statuses {
		if !slices.Contains(types.AllEventStatuses, status) {
			invalid = append(invalid, string(status))
		}
	}
	if len(invalid) < 1 {
		return nil
	}
	return core.InvalidEventStatusError(pc, invalid)
}

func (c *EventConfiguration) Normalize(ctx context.Context, pc core.ObjectPathContext, resolver core.NormalizeResolver) core.NormalizeErrors {
	errs := core.NormalizeErrors{}
	for _, blockName := range c.BlockNames {
		block, err := resolver.ResolveBlock(ctx, types.ConnectionTarget{BlockName: blockName})
		if err != nil {
			errs = append(errs, core.NormalizeError{
				ObjectPathContext: pc.SubKey("blocks", blockName),
				ErrorMessage:      fmt.Sprintf("Block must exist to subscribe to event, but it failed to resolve: %s", err),
			})
		} else {
			c.Blocks = append(c.Blocks, block)
		}
	}
	return errs
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

func (c *EventConfiguration) ToEnvEvent(name string, env types.Environment) types.EnvEvent {
	var channels map[types.IntegrationTool]types.ChannelData
	if c.Targets != nil {
		channels = c.Targets.Channels()
	}
	return types.EnvEvent{
		OrgName:  env.OrgName,
		StackId:  env.StackId,
		EnvId:    env.Id,
		Name:     name,
		Actions:  c.Actions,
		Blocks:   c.BlockIds(),
		Statuses: c.Statuses,
		Channels: channels,
	}
}

func (c *EventConfiguration) OverrideEnvEvent(event types.EnvEvent) types.EnvEvent {
	if c.Actions != nil {
		event.Actions = c.Actions
	}
	if c.Statuses != nil {
		event.Statuses = c.Statuses
	}
	if c.Blocks != nil {
		event.Blocks = c.BlockIds()
	}
	if c.Targets != nil {
		event.Channels = c.Targets.Channels()
	}
	return event
}

func (c *EventConfiguration) BlockIds() []int64 {
	if c.Blocks == nil {
		return nil
	}
	blocks := make([]int64, 0)
	for _, cur := range c.Blocks {
		blocks = append(blocks, cur.Id)
	}
	return blocks
}
