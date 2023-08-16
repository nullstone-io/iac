package services

import (
	"context"
	"gopkg.in/nullstone-io/go-api-client.v0"
	"gopkg.in/nullstone-io/go-api-client.v0/auth"
)

func NewRunnerApiHub(config ApiHubConfig, store auth.RunnerKeyStore) ApiHub {
	authServer := &auth.NsAuthServer{Address: config.AuthServerAddr}

	ats := &auth.RunnerAccessTokenSource{
		AuthServer:     authServer,
		RunnerKeyStore: store,
	}

	return ApiHub{
		MainFactory: &api.ClientFactory{
			BaseAddress:       config.MainAddr,
			AccessTokenSource: ats,
			IsTraceEnabled:    config.IsTraceEnabled,
		},
		FurionFactory: &api.ClientFactory{
			BaseAddress:       config.FurionAddr,
			AccessTokenSource: ats,
			IsTraceEnabled:    config.IsTraceEnabled,
		},
		EnigmaFactory: &api.ClientFactory{
			BaseAddress:       config.EnigmaAddr,
			AccessTokenSource: ats,
			IsTraceEnabled:    config.IsTraceEnabled,
		},
	}
}

func NewApiHub(config ApiHubConfig, source auth.AccessTokenSource) ApiHub {
	return ApiHub{
		MainFactory: &api.ClientFactory{
			BaseAddress:       config.MainAddr,
			AccessTokenSource: source,
			IsTraceEnabled:    config.IsTraceEnabled,
		},
		FurionFactory: &api.ClientFactory{
			BaseAddress:       config.FurionAddr,
			AccessTokenSource: source,
			IsTraceEnabled:    config.IsTraceEnabled,
		},
		EnigmaFactory: &api.ClientFactory{
			BaseAddress:       config.EnigmaAddr,
			AccessTokenSource: source,
			IsTraceEnabled:    config.IsTraceEnabled,
		},
	}
}

type ApiHub struct {
	MainFactory   *api.ClientFactory
	FurionFactory *api.ClientFactory
	EnigmaFactory *api.ClientFactory
}

type apiHubContextKey struct{}

func ApiHubFromContext(ctx context.Context) ApiHub {
	if val, ok := ctx.Value(apiHubContextKey{}).(ApiHub); ok {
		return val
	}
	return ApiHub{}
}
func ContextWithApiHub(ctx context.Context, apiHub ApiHub) context.Context {
	return context.WithValue(ctx, apiHubContextKey{}, apiHub)
}

func (h ApiHub) Client(orgName string) *api.Client {
	return h.MainFactory.Client(orgName)
}

func (h ApiHub) Furion(orgName string) *api.Client {
	return h.FurionFactory.Client(orgName)
}

func (h ApiHub) Enigma(orgName string) *api.Client {
	return h.EnigmaFactory.Client(orgName)
}

func (h ApiHub) GetClaims(orgName string) (*auth.Claims, error) {
	return auth.GetCustomClaims(h.MainFactory.AccessTokenSource, orgName)
}
