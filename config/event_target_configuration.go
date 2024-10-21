package config

import (
	"fmt"
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
		events[target] = eventTargetConfigFromYaml(value)
	}
	return events
}

func (s EventTargetConfigurations) Validate(ic core.IacContext, pc core.ObjectPathContext) core.ValidateErrors {
	errs := make(core.ValidateErrors, 0)
	for target, cur := range s {
		errs = append(errs, cur.Validate(ic, pc.SubKey("targets", target))...)
	}
	return errs
}

type EventTargetConfiguration struct {
	Target string `json:"target"`
}

func eventTargetConfigFromYaml(target string, value yaml.EventTargetConfiguration) *EventTargetConfiguration {
	return &EventTargetConfiguration{
		Target: target,
	}
}

func (c *EventTargetConfiguration) Validate(ic core.IacContext, pc core.ObjectPathContext) core.ValidateErrors {
	errs := core.ValidateErrors{}
	if _, ok := AllEventTargets[c.Target]; !ok {
		errs = append(errs, core.InvalidEventTargetError(pc, c.Target))
	} else {
		if validatorFn, ok := eventTypeDataValidators[c.Type]; ok {
			errs = append(errs, validatorFn(pc, c.Data)...)
		}
	}
	return errs
}

type EventTypeDataValidatorFunc func(core.ObjectPathContext, map[string]any) core.ValidateErrors

var (
	eventTypeDataValidators = map[EventType]EventTypeDataValidatorFunc{
		EventTypeSlackNotification: func(pc core.ObjectPathContext, data map[string]any) core.ValidateErrors {
			errs := core.ValidateErrors{}
			if val, ok := data["workspace"]; !ok || val == nil {
				errs = append(errs, core.MissingRequiredEventData(pc.SubField("data"), "workspace"))
			}
			if val, ok := data["channels"]; !ok || val == nil {
				errs = append(errs, core.MissingRequiredEventData(pc.SubField("data"), "channels"))
			} else if vals, ok := val.([]any); !ok || len(vals) == 0 {
				errs = append(errs, core.ValidateError{
					ObjectPathContext: pc.SubKey("data", "channels"),
					ErrorMessage:      fmt.Sprintf("A %s event requires at least one channel", EventTypeSlackNotification),
				})
			}
			return errs
		},
	}
)
