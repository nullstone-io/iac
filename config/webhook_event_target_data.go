package config

import (
	"context"
	"fmt"
	"net/url"

	"github.com/nullstone-io/iac/core"
	"github.com/nullstone-io/iac/yaml"
)

type WebhookEventTargetData struct {
	Url string `json:"url"`
}

func webhookEventTargetDataFromYaml(yml *yaml.EventTargetWebhookConfiguration) *WebhookEventTargetData {
	if yml == nil {
		return nil
	}
	return &WebhookEventTargetData{Url: yml.Url}
}

func (d *WebhookEventTargetData) Resolve(ctx context.Context, resolver core.EventChannelResolver, ic core.IacContext, pc core.ObjectPathContext) core.ResolveErrors {
	return nil
}

func (d *WebhookEventTargetData) Validate(ic core.IacContext, pc core.ObjectPathContext) core.ValidateErrors {
	if d == nil {
		return nil
	}
	errs := core.ValidateErrors{}
	if d.Url == "" {
		errs = append(errs, core.ValidateError{
			ObjectPathContext: pc,
			ErrorMessage:      fmt.Sprintf("When specifying `webhook`, `url` is required"),
		})
	} else {
		parsed, err := url.Parse(d.Url)
		if err != nil {
			errs = append(errs, core.ValidateError{
				ObjectPathContext: pc.SubField("url"),
				ErrorMessage:      fmt.Sprintf("Invalid webhook URL: %s", err.Error()),
			})
		} else if parsed.Scheme == "" {
			errs = append(errs, core.ValidateError{
				ObjectPathContext: pc.SubField("url"),
				ErrorMessage:      fmt.Sprintf("Invalid webhook URL"),
			})
		}
	}
	return errs
}

func (d *WebhookEventTargetData) ChannelData() map[string]any {
	connections := make([]map[string]any, 0)
	if d.Url != "" {
		connections = append(connections, map[string]any{"url": d.Url})
	}
	return map[string]any{"connections": connections}
}
