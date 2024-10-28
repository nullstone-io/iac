package config

import (
	"context"
	"github.com/nullstone-io/iac/core"
	"github.com/nullstone-io/iac/yaml"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

var (
	AllEventTargets = map[types.EventTarget]bool{}
)

type EventTargetConfigurations map[string]*EventTargetConfiguration

func convertEventTargetConfigurations(parsed map[string]yaml.EventTargetConfiguration) EventTargetConfigurations {
	events := EventTargetConfigurations{}
	for target, value := range parsed {
		events[target] = eventTargetConfigFromYaml(target, value)
	}
	return events
}

func (s EventTargetConfigurations) Resolve(ctx context.Context, resolver core.EventChannelResolver, ic core.IacContext, pc core.ObjectPathContext) core.ResolveErrors {
	errs := core.ResolveErrors{}
	for name, cur := range s {
		errs = append(errs, cur.Resolve(ctx, resolver, ic, pc.SubKey("targets", name))...)
	}
	return errs
}

func (s EventTargetConfigurations) Validate(ic core.IacContext, pc core.ObjectPathContext) core.ValidateErrors {
	errs := make(core.ValidateErrors, 0)
	for target, cur := range s {
		errs = append(errs, cur.Validate(ic, pc.SubKey("targets", target))...)
	}
	return errs
}

func (s EventTargetConfigurations) Channels() map[types.IntegrationTool]types.ChannelData {
	result := map[types.IntegrationTool]types.ChannelData{}
	for _, cur := range s {
		result[types.IntegrationTool(cur.Target)] = cur.ChannelData()
	}
	return result
}

type EventTargetConfiguration struct {
	Target string `json:"target"`

	SlackData *SlackEventTargetData `json:"slackData"`
}

func eventTargetConfigFromYaml(target string, value yaml.EventTargetConfiguration) *EventTargetConfiguration {
	return &EventTargetConfiguration{
		Target:    target,
		SlackData: slackEventTargetDataFromYaml(value.Slack),
	}
}

func (c *EventTargetConfiguration) ChannelData() map[string]any {
	if c.SlackData != nil {
		return c.SlackData.ChannelData()
	}
	return nil
}

func (c *EventTargetConfiguration) Resolve(ctx context.Context, resolver core.EventChannelResolver, ic core.IacContext, pc core.ObjectPathContext) core.ResolveErrors {
	errs := core.ResolveErrors{}
	if c.SlackData != nil {
		errs = append(errs, c.SlackData.Resolve(ctx, resolver, ic, pc.SubField("slack"))...)
	}
	return errs
}

func (c *EventTargetConfiguration) Validate(ic core.IacContext, pc core.ObjectPathContext) core.ValidateErrors {
	errs := core.ValidateErrors{}
	eventTarget := types.EventTarget(c.Target)
	if _, ok := AllEventTargets[eventTarget]; !ok {
		errs = append(errs, core.InvalidEventTargetError(pc, c.Target))
	} else {
		errs = append(errs, c.SlackData.Validate(ic, pc.SubField("slack"))...)
	}
	return errs
}
