package core

import (
	"context"
)

var (
	_ EventChannelResolver = &StaticEventChannelResolver{}
)

type StaticEventChannelResolver struct {
	ChannelsByTool map[string][]map[string]any
}

func (s StaticEventChannelResolver) ListChannels(ctx context.Context, tool string) ([]map[string]any, error) {
	return s.ChannelsByTool[tool], nil
}
