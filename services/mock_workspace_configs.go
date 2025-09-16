package services

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

type MockWorkspaceConfigEntry struct {
	StackId int64
	BlockId int64
	EnvId   int64
	Config  types.WorkspaceConfig
}

func MockLatestWorkspaceConfigs(router *mux.Router, entries []MockWorkspaceConfigEntry) {
	find := func(stackId, blockId, envId int64) *types.WorkspaceConfig {
		for _, entry := range entries {
			if entry.StackId == stackId && entry.BlockId == blockId && entry.EnvId == envId {
				return &entry.Config
			}
		}
		return nil
	}

	router.Path("/orgs/{orgName}/stacks/{stackId}/blocks/{blockId}/envs/{envId}/configs/latest").
		Methods(http.MethodGet).
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			stackId, err := strconv.ParseInt(vars["stackId"], 10, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			blockId, err := strconv.ParseInt(vars["blockId"], 10, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			envId, err := strconv.ParseInt(vars["envId"], 10, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			config := find(stackId, blockId, envId)
			if config == nil {
				http.NotFound(w, r)
				return
			}
			raw, _ := json.Marshal(config)
			w.Write(raw)
		})
}
