package services

import (
	"gopkg.in/nullstone-io/go-api-client.v0"
	"gopkg.in/nullstone-io/go-api-client.v0/auth"
)

func NewApiHub(config ApiHubConfig, source auth.AccessTokenSource) ApiHub {
	return ApiHub{
		MainFactory: &api.ClientFactory{
			BaseAddress:       config.MainAddr,
			AccessTokenSource: source,
			IsTraceEnabled:    config.IsTraceEnabled,
		},
	}
}

type ApiHub struct {
	MainFactory *api.ClientFactory
}

func (h ApiHub) Client(orgName string) *api.Client {
	return h.MainFactory.Client(orgName)
}
