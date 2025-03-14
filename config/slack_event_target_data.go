package config

import (
	"context"
	"fmt"
	"github.com/nullstone-io/iac/core"
	"github.com/nullstone-io/iac/yaml"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

type SlackEventTargetData struct {
	// Channels refers to the Slack channels by name
	Channels []string `json:"channels"`

	// ChannelIds are loaded during Resolve
	ChannelIds map[string]string `json:"channelIds"`
}

func slackEventTargetDataFromYaml(yml *yaml.EventTargetSlackConfiguration) *SlackEventTargetData {
	if yml == nil {
		return nil
	}
	return &SlackEventTargetData{Channels: yml.Channels}
}

func (d *SlackEventTargetData) Resolve(ctx context.Context, resolver core.EventChannelResolver, ic core.IacContext, pc core.ObjectPathContext) core.ResolveErrors {
	channels, err := resolver.ListChannels(ctx, string(types.IntegrationToolSlack))
	if err != nil {
		return core.ResolveErrors{core.ToolChannelLookupFailedError(pc, string(types.IntegrationToolSlack), err)}
	}

	allChannelIds := map[string]string{}
	for _, rawData := range channels {
		var id string
		var name string
		if val, ok := rawData["id"]; ok {
			id, _ = val.(string)
		}
		if val, ok := rawData["name"]; ok {
			name, _ = val.(string)
		}
		if name != "" {
			allChannelIds[name] = id
		}
	}

	errs := core.ResolveErrors{}
	d.ChannelIds = map[string]string{}
	for _, channelName := range d.Channels {
		if id, ok := allChannelIds[channelName]; ok {
			d.ChannelIds[channelName] = id
		} else {
			errs = append(errs, core.ResolveError{
				ObjectPathContext: pc,
				ErrorMessage:      fmt.Sprintf("Slack channel %q could not be found in the connected workspace", channelName),
			})
		}
	}
	return errs
}

func (d *SlackEventTargetData) Validate(ic core.IacContext, pc core.ObjectPathContext) core.ValidateErrors {
	if d == nil {
		return nil
	}
	return nil
}

func (d *SlackEventTargetData) ChannelData() map[string]any {
	connections := make([]map[string]any, 0)
	for _, channelId := range d.ChannelIds {
		connections = append(connections, map[string]any{"channel_id": channelId})
	}
	return map[string]any{"connections": connections}
}
