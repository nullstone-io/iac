package services

import "gopkg.in/nullstone-io/go-api-client.v0"

type ApiHubConfig struct {
	MainAddr       string
	FurionAddr     string
	EnigmaAddr     string
	AuthServerAddr string
	IsTraceEnabled bool
}

func DefaultApiHubConfig() ApiHubConfig {
	cfg := api.DefaultConfig()
	return ApiHubConfig{
		MainAddr:       cfg.BaseAddress,
		FurionAddr:     FurionAddress.Get(),
		EnigmaAddr:     EnigmaAddress.Get(),
		AuthServerAddr: VoidAddress.Get(),
		IsTraceEnabled: cfg.IsTraceEnabled,
	}
}
