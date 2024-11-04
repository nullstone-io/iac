package core

import (
	"context"
	"fmt"
	"gopkg.in/nullstone-io/go-api-client.v0"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
	"sync"
)

var (
	_ EventChannelResolver = &ApiEventChannelResolver{}
)

type ApiEventChannelResolver struct {
	ApiClient *api.Client

	once         sync.Once
	loadErr      error
	integrations []types.Integration
	cache        map[string][]map[string]any
}

func (a *ApiEventChannelResolver) ListChannels(ctx context.Context, tool string) ([]map[string]any, error) {
	a.once.Do(func() {
		a.load(ctx)
	})
	if a.loadErr != nil {
		return nil, a.loadErr
	}

	byTool, ok := a.cache[tool]
	if ok {
		return byTool, nil
	}

	integration := a.findIntegration(tool)
	if integration == nil {
		return nil, fmt.Errorf("No %q integration found in %s organization.", tool, a.ApiClient.Config.OrgName)
	}
	status, err := a.ApiClient.Integrations().GetStatus(ctx, integration.Id)
	if err != nil {
		return nil, err
	} else if status == nil {
		return nil, fmt.Errorf("Unable to find channels for %q in %s organization", tool, a.ApiClient.Config.OrgName)
	}

	byTool = []map[string]any{}
	for _, cur := range status.SlackChannels {
		byTool = append(byTool, map[string]any{
			"id":              cur.ID,
			"name":            cur.Name,
			"is_private":      cur.IsPrivate,
			"is_im":           cur.IsIM,
			"context_team_id": cur.ContextTeamId,
		})
	}
	return byTool, nil
}

func (a *ApiEventChannelResolver) load(ctx context.Context) {
	a.cache = map[string][]map[string]any{}
	a.integrations, a.loadErr = a.ApiClient.Integrations().List(ctx)
	if a.integrations == nil {
		a.loadErr = fmt.Errorf("Nullstone integrations could not be found")
	}
}

func (a *ApiEventChannelResolver) findIntegration(tool string) *types.Integration {
	for _, cur := range a.integrations {
		if cur.Tool == types.IntegrationTool(tool) {
			return &cur
		}
	}
	return nil
}
