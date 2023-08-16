package oracle

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
	"net/http"
)

func MockGetModuleVersions(router *mux.Router, modules ...*types.Module) {
	findModule := func(orgName, name string) *types.Module {
		for _, m := range modules {
			if m.OrgName == orgName && m.Name == name {
				return m
			}
		}
		return nil
	}

	router.Path("/orgs/{orgName}/modules/{moduleName}").
		Methods(http.MethodGet).
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			block := findModule(vars["orgName"], vars["moduleName"])
			if block == nil {
				http.NotFound(w, r)
			} else {
				raw, _ := json.Marshal(block)
				w.Write(raw)
			}
		})

	findModuleVersion := func(orgName, name, version string) *types.ModuleVersion {
		for _, m := range modules {
			if m.OrgName == orgName && m.Name == name {
				for _, v := range m.Versions {
					if version == "latest" || v.Version == version {
						return &v
					}
				}
			}
		}
		return nil
	}

	router.Path("/orgs/{orgName}/modules/{moduleName}/versions/{version}").
		Methods(http.MethodGet).
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			block := findModuleVersion(vars["orgName"], vars["moduleName"], vars["version"])
			if block == nil {
				http.NotFound(w, r)
			} else {
				raw, _ := json.Marshal(block)
				w.Write(raw)
			}
		})
}
