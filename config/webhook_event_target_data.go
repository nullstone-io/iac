package config

import (
	"context"
	"fmt"
	"net/url"

	"github.com/nullstone-io/iac/core"
	"github.com/nullstone-io/iac/yaml"
)

type WebhookEventTargetData struct {
	Urls []string `json:"urls"`
}

func webhookEventTargetDataFromYaml(yml *yaml.EventTargetWebhookConfiguration) *WebhookEventTargetData {
	if yml == nil {
		return nil
	}
	return &WebhookEventTargetData{Urls: yml.Urls}
}

func (d *WebhookEventTargetData) Resolve(ctx context.Context, resolver core.EventChannelResolver, ic core.IacContext, pc core.ObjectPathContext) core.ResolveErrors {
	return nil
}

func (d *WebhookEventTargetData) Validate(ic core.IacContext, pc core.ObjectPathContext) core.ValidateErrors {
	if d == nil {
		return nil
	}
	errs := core.ValidateErrors{}
	if len(d.Urls) == 0 {
		errs = append(errs, core.ValidateError{
			ObjectPathContext: pc,
			ErrorMessage:      fmt.Sprintf("When specifying `webhook`, it must have at least one url in 'urls'"),
		})
	}
	for i, val := range d.Urls {
		parsed, err := url.Parse(val)
		if err != nil {
			errs = append(errs, core.ValidateError{
				ObjectPathContext: pc.SubIndex("urls", i),
				ErrorMessage:      fmt.Sprintf("Invalid webhook URL: %s", err.Error()),
			})
		} else if parsed.Scheme == "" {
			errs = append(errs, core.ValidateError{
				ObjectPathContext: pc.SubIndex("urls", i),
				ErrorMessage:      fmt.Sprintf("Invalid webhook URL"),
			})
		}
	}
	return errs
}

func (d *WebhookEventTargetData) ChannelData() map[string]any {
	connections := make([]map[string]any, 0)
	for _, val := range d.Urls {
		connections = append(connections, map[string]any{"incoming_webhook": map[string]any{"url": val}})
	}
	return map[string]any{"connections": connections}
}
