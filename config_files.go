package iac

import (
	"github.com/nullstone-io/iac/config"
	"github.com/nullstone-io/iac/core"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

type ConfigFiles struct {
	// RepoUrl refers to the vcs repository containing these files (e.g. https://github.com/acme/repo)
	RepoUrl string `json:"repoUrl"`

	// RepoName refers to the name of the vcs repository (e.g. acme/repo)
	RepoName string `json:"repoName"`

	// Config contains `.nullstone/config.yml` validated and normalized as Nullstone objects
	// This was added to the state since TemporalIacSync.IacSync Config is intentionally redacted from json
	Config *config.EnvConfiguration `json:"config"`

	// Overrides contains `.nullstone/<env|previews>.yml` validated and normalized as Nullstone objects
	// This was added to the state since TemporalIacSync.IacSync Overrides is intentionally redacted from json
	Overrides map[string]*config.EnvConfiguration `json:"overrides"`
}

func (r ConfigFiles) BlockNames(env types.Environment) map[string]bool {
	blockNames := map[string]bool{}
	if r.Config != nil {
		blockNames = r.Config.BlockNames()
	}
	if cur := r.GetOverrides(env); cur != nil {
		for k := range cur.BlockNames() {
			blockNames[k] = true
		}
	}
	return blockNames
}

func (r ConfigFiles) GetOverrides(env types.Environment) *config.EnvConfiguration {
	if env.Type == types.EnvTypePreview {
		ec, _ := r.Overrides["previews"]
		return ec
	}
	ec, _ := r.Overrides[env.Name]
	return ec
}

func (r ConfigFiles) NewIacFinder(env types.Environment) core.IacFinder {
	return config.NewIacFinder(r.Config, r.GetOverrides(env), env.StackId, env.Id)
}
