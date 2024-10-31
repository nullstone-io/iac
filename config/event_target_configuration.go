package config

import (
	"context"
	"github.com/nullstone-io/iac/core"
	"github.com/nullstone-io/iac/yaml"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

var (
	AllEventTargets = map[types.EventTarget]bool{
		types.EventTargetSlack:          true,
		types.EventTargetMicrosoftTeams: true,
		types.EventTargetDiscord:        true,
		types.EventTargetWhatsapp:       true,
		types.EventTargetWebhook:        true,
		types.EventTargetTask:           true,
	}
)

type EventTargetConfigurations map[string]*EventTargetConfiguration

func convertEventTargetConfigurations(parsed yaml.EventTargetConfiguration) EventTargetConfigurations {
	events := EventTargetConfigurations{}
	if parsed.Slack != nil {
		events[string(types.EventTargetSlack)] = &EventTargetConfiguration{
			Target:    string(types.EventTargetSlack),
			SlackData: slackEventTargetDataFromYaml(parsed.Slack),
		}
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

func (c *EventTargetConfiguration) ChannelData() map[string]any {
	if c.SlackData != nil {
		return c.SlackData.ChannelData()
	}
	return nil
}

func (c *EventTargetConfiguration) Resolve(ctx context.Context, resolver core.EventChannelResolver, ic core.IacContext, pc core.ObjectPathContext) core.ResolveErrors {
	errs := core.ResolveErrors{}
	if c.SlackData != nil {
		errs = append(errs, c.SlackData.Resolve(ctx, resolver, ic, pc)...)
	}
	return errs
}

func (c *EventTargetConfiguration) Validate(ic core.IacContext, pc core.ObjectPathContext) core.ValidateErrors {
	errs := core.ValidateErrors{}
	eventTarget := types.EventTarget(c.Target)
	if _, ok := AllEventTargets[eventTarget]; !ok {
		errs = append(errs, core.InvalidEventTargetError(pc, c.Target))
	} else {
		errs = append(errs, c.SlackData.Validate(ic, pc)...)
	}
	return errs
}
