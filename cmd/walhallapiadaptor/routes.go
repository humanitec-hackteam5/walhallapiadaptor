package main

import "github.com/gorilla/mux"

func (s *server) setupRoutes() {
	r := mux.NewRouter()
	r.Methods("GET").Path("/orgs").HandlerFunc(s.listOrgs())
	r.Methods("GET").Path("/orgs/{orgId}/modules").HandlerFunc(s.listModules())
	r.Methods("POST").Path("/orgs/{orgId}/modules/refresh").HandlerFunc(s.refreshModules())
	r.Methods("GET").Path("/orgs/{orgId}/modules/refresh").HandlerFunc(s.getRefreshModulesStatus())
	//r.Methods("GET").Path("/orgs/modules/{moduleName}").HandlerFunc(s.getModule())
	//r.Methods("GET").Path("/orgs/modules/{moduleName}/build").HandlerFunc(s.listModuleBuilds())
	//r.Methods("GET").Path("/orgs/modules/{moduleName}/build/").HandlerFunc(s.getModuleBuild())
	s.router = r
}
