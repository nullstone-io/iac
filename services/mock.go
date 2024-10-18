package services

import (
	"gopkg.in/nullstone-io/go-api-client.v0/auth"
	"gopkg.in/nullstone-io/go-api-client.v0/trace"
	"net/http"
	"net/http/httptest"
	"testing"
)

func MockApiHub(t *testing.T, handler http.Handler) ApiHub {
	server := httptest.NewServer(handler)
	cfg := ApiHubConfig{
		MainAddr:       server.URL,
		IsTraceEnabled: trace.IsEnabled(),
	}
	apiHub := NewApiHub(cfg, auth.RawAccessTokenSource{AccessToken: "mock-user"})
	t.Cleanup(server.Close)
	return apiHub
}
