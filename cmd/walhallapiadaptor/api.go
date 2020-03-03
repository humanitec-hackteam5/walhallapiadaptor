package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Module struct {
	Name   string        `json:"name"`
	Source string        `json:"source"`
	Builds []ModuleBuild `json:"builds"`
}
type ModuleBuild struct {
	Image  string   `json:"image"`
	Commit string   `json:"commit"`
	Branch string   `json:"branch"`
	Tags   []string `json:"tags"`
}

// listOrgs returns a handler which returns a list of all the orgs the user is a member of
//
func (s *server) listOrgs() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		walhall, err := s.newWalhall(r.Header.Get("authorization"))
		if err != nil {
			w.WriteHeader(403)
			fmt.Fprint(w, `"Unable to parse JWT"`)
			return
		}
		orgs, err := walhall.ListOrgs()
		if err != nil {
			log.Printf("list apps: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		orgNames := make([]string, len(orgs))
		i := 0
		for k := range orgs {
			orgNames[i] = k
			i++
		}
		encoder := json.NewEncoder(w)
		err = encoder.Encode(orgNames)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}
	}
}

// listModules returns a handler which returns a list of all the modules available to the user in an org
//
func (s *server) listModules() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		walhall, err := s.newWalhall(r.Header.Get("authorization"))
		if err != nil {
			w.WriteHeader(403)
			fmt.Fprint(w, `"Unable to parse JWT"`)
			return
		}
		walhallModules, err := walhall.ListModules(params["orgId"])
		if err != nil {
			log.Printf("list modules: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		modules := make([]Module, len(walhallModules))
		for iM, module := range walhallModules {
			builds := make([]ModuleBuild, len(module.Versions))
			for iV, version := range module.Versions {
				builds[iV] = ModuleBuild{
					Image:  fmt.Sprintf("%s/%s/%s:%s", s.registryName, params["orgId"], module.Image, version.Version),
					Commit: "UNKNOWN",
					Branch: "UNKNOWN",
					Tags:   []string{version.Version},
				}
			}
			modules[iM] = Module{
				Name:   module.Name,
				Source: "Github",
				Builds: builds,
			}
		}

		encoder := json.NewEncoder(w)
		err = encoder.Encode(modules)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}
	}
}

// refreshModules returns a handler which forces Walhall to refresh the module list on the BE
//
func (s *server) refreshModules() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		walhall, err := s.newWalhall(r.Header.Get("authorization"))
		if err != nil {
			w.WriteHeader(403)
			fmt.Fprint(w, `"Unable to parse JWT"`)
			return
		}
		status, err := walhall.RefreshModules(params["orgId"])
		if err != nil {
			log.Printf("list modules: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		encoder := json.NewEncoder(w)
		err = encoder.Encode(status)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}
	}
}

// getRefreshModulesStatus returns a handler which checks the status of the refresh of modules
//
func (s *server) getRefreshModulesStatus() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		walhall, err := s.newWalhall(r.Header.Get("authorization"))
		if err != nil {
			w.WriteHeader(403)
			fmt.Fprint(w, `"Unable to parse JWT"`)
			return
		}
		status, err := walhall.GetRefreshModulesStatus(params["orgId"])
		if err != nil {
			log.Printf("list modules: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		encoder := json.NewEncoder(w)
		err = encoder.Encode(status)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}
	}
}
