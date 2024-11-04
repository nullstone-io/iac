package core

import (
	"context"
	"fmt"
)

var (
	_ EventChannelResolver = &StaticEventChannelResolver{}
)

type StaticEventChannelResolver struct {
	ChannelsByTool map[string][]map[string]any
}

func (s StaticEventChannelResolver) ListChannels(ctx context.Context, tool string) ([]map[string]any, error) {
	if s.ChannelsByTool == nil {
		return nil, ErrEventChannelsNotInitialized
	}
	byTool, ok := s.ChannelsByTool[tool]
	if !ok {
		return nil, EventChannelsByToolsNotInitializedError{Tool: tool}
	}
	return byTool, nil
}

var (
	_ error = EventChannelsByToolsNotInitializedError{}
)

type EventChannelsByToolsNotInitializedError struct {
	Tool string
}

func (e EventChannelsByToolsNotInitializedError) Error() string {
	return fmt.Sprintf("event channels have not been initialized for %s", e.Tool)
}
