package core

import (
	"context"
	"errors"
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
		return nil, NoIntegrationFoundError{OrgName: a.ApiClient.Config.OrgName, Tool: tool}
	}
	status, err := a.ApiClient.Integrations().GetStatus(ctx, integration.Id)
	if err != nil {
		return nil, err
	} else if status == nil {
		return nil, fmt.Errorf("Unable to find channels for %q in %s organization", tool, a.ApiClient.Config.OrgName)
	}

	if !status.IsConnected {
		return nil, IntegrationDisconnectedError{OrgName: a.ApiClient.Config.OrgName, Tool: tool}
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
	a.cache[tool] = byTool
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

var (
	_ error = NoIntegrationFoundError{}
	_ error = IntegrationDisconnectedError{}
)

type NoIntegrationFoundError struct {
	OrgName string
	Tool    string
}

func (e NoIntegrationFoundError) Error() string {
	return fmt.Sprintf("No %q integration found in %s organization.", e.Tool, e.OrgName)
}

func IsNoIntegrationFoundError(err error) bool {
	var nife NoIntegrationFoundError
	return errors.As(err, &nife)
}

type IntegrationDisconnectedError struct {
	OrgName string
	Tool    string
}

func (e IntegrationDisconnectedError) Error() string {
	return fmt.Sprintf("%q integration is disconnected in %s organization.", e.Tool, e.OrgName)
}

func IsIntegrationDisconnectedError(err error) bool {
	var ide IntegrationDisconnectedError
	return errors.As(err, &ide)
}
