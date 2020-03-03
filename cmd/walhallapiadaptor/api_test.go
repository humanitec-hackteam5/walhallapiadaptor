package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/matryer/is"
	"humanitec.io/walhallapiadaptor/internal/walhallapi"
)

// NOTE: *_mock.go files are generated via the following commands:
// $ mockgen -source=../../internal/walhallapi/types.go -destination=walhallapier_mock.go -package=main WalhallAPIer

type mocks struct {
	walhall  walhallapi.WalhallAPIer
	registry string
}

func ExecuteRequest(mocks mocks, method, url string, body io.Reader, t *testing.T) *httptest.ResponseRecorder {
	server := server{
		newWalhall: func(jwt string) (walhallapi.WalhallAPIer, error) {
			return mocks.walhall, nil
		},
		registryName: mocks.registry,
	}
	server.setupRoutes()

	var req *http.Request
	var err error
	if body == nil {
		req, err = http.NewRequest(method, url, nil)
	} else {
		req, err = http.NewRequest(method, url, body)
	}
	if err != nil {
		t.Errorf("creating request: %v", err)
	}

	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)
	return w
}

func containsSameValues(aIn, bIn interface{}) bool {
	aType := reflect.ValueOf(aIn)
	bType := reflect.ValueOf(bIn)
	if aType.Kind() != reflect.Slice || bType.Kind() != reflect.Slice {
		return false
	}

	if aType.Len() != bType.Len() {
		log.Printf("%v != %v\n", aIn, bIn)
		return false
	}

	matchA := make([]int, aType.Len())
	matchB := make([]int, aType.Len())

	for i := 0; i < aType.Len(); i++ {
		for j := 0; j < aType.Len(); j++ {
			// Count number of duplicates in a
			if reflect.DeepEqual(aType.Index(i).Interface(), aType.Index(j).Interface()) {
				matchA[i]++
			}
			// Count number of times a[i] matches in b
			if reflect.DeepEqual(aType.Index(i).Interface(), bType.Index(j).Interface()) {
				matchB[i]++
			}
		}
	}

	for i := range matchA {
		if matchA[i] != matchB[i] {
			return false
		}
	}

	return true
}

func TestListOrgs(t *testing.T) {
	is := is.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockWalhallAPIer(ctrl)

	expected := []string{"org-one", "org-two"}
	m.
		EXPECT().
		ListOrgs().
		Return(map[string]string{
			"org-one": "ORGID01",
			"org-two": "ORGID02",
		}, nil).
		Times(1)

	resp := ExecuteRequest(mocks{walhall: m}, http.MethodGet, "/orgs", nil, t)

	var actual []string
	json.Unmarshal(resp.Body.Bytes(), &actual)
	is.True(containsSameValues(expected, actual))
}

func TestListModules(t *testing.T) {
	is := is.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedModules := []Module{
		Module{
			Name:   "test-module-one",
			Source: "Github",
			Builds: []ModuleBuild{
				ModuleBuild{
					Image:  "registry.walhall.io/org-one/test-module-one:VERSION_ONE",
					Commit: "UNKNOWN",
					Branch: "UNKNOWN",
					Tags:   []string{"VERSION_ONE"},
				},
				ModuleBuild{
					Image:  "registry.walhall.io/org-one/test-module-one:VERSION_TWO",
					Commit: "UNKNOWN",
					Branch: "UNKNOWN",
					Tags:   []string{"VERSION_TWO"},
				},
			},
		},
		Module{
			Name:   "test-module-two",
			Source: "Github",
			Builds: []ModuleBuild{
				ModuleBuild{
					Image:  "registry.walhall.io/org-one/test-module-two:VERSION_ONE",
					Commit: "UNKNOWN",
					Branch: "UNKNOWN",
					Tags:   []string{"VERSION_ONE"},
				},
				ModuleBuild{
					Image:  "registry.walhall.io/org-one/test-module-two:VERSION_TWO",
					Commit: "UNKNOWN",
					Branch: "UNKNOWN",
					Tags:   []string{"VERSION_TWO"},
				},
			},
		},
	}
	m := NewMockWalhallAPIer(ctrl)

	m.
		EXPECT().
		ListModules("org-one").
		Return([]walhallapi.Module{
			walhallapi.Module{
				Name:  "test-module-one",
				UUID:  "5a6b1ac7-6bd5-4b07-827b-047a94ccc91a",
				Repo:  "org-one/test-module-one",
				Image: "test-module-one",
				Versions: []walhallapi.ModuleVersion{
					walhallapi.ModuleVersion{
						ID:      1001,
						UUID:    "59304d84-d503-44cf-a171-00367d8bacb4",
						Version: "VERSION_ONE",
					},
					walhallapi.ModuleVersion{
						ID:      1002,
						UUID:    "44bc53dd-142a-41a4-9d29-896f5fb3f0d0",
						Version: "VERSION_TWO",
					},
				},
			},
			walhallapi.Module{
				Name:  "test-module-two",
				UUID:  "2273b8c5-5704-433f-91ae-8471de6b4f5f",
				Repo:  "org-one/test-module-two",
				Image: "test-module-two",
				Versions: []walhallapi.ModuleVersion{
					walhallapi.ModuleVersion{
						ID:      2001,
						UUID:    "9a7bf0e1-2386-4da2-80c4-3e81df1ebac4",
						Version: "VERSION_ONE",
					},
					walhallapi.ModuleVersion{
						ID:      2002,
						UUID:    "d4705bb7-abc5-4efd-b7e3-bbd041d4aebd",
						Version: "VERSION_TWO",
					},
				},
			},
		}, nil).
		Times(1)

	resp := ExecuteRequest(mocks{walhall: m, registry: "registry.walhall.io"}, http.MethodGet, "/orgs/org-one/modules", nil, t)

	var actual []Module
	json.Unmarshal(resp.Body.Bytes(), &actual)
	is.Equal(actual, expectedModules)
}

func TestRefreshModules(t *testing.T) {
	is := is.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockWalhallAPIer(ctrl)

	m.
		EXPECT().
		RefreshModules("org-one").
		Return("running", nil).
		Times(1)

	resp := ExecuteRequest(mocks{walhall: m}, http.MethodPost, "/orgs/org-one/modules/refresh", nil, t)

	var actual string
	json.Unmarshal(resp.Body.Bytes(), &actual)
	is.Equal(actual, "running")
}

func TestGetRefreshModulesStatus(t *testing.T) {
	is := is.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockWalhallAPIer(ctrl)

	m.
		EXPECT().
		GetRefreshModulesStatus("org-one").
		Return("success", nil).
		Times(1)

	resp := ExecuteRequest(mocks{walhall: m}, http.MethodGet, "/orgs/org-one/modules/refresh", nil, t)

	var actual string
	json.Unmarshal(resp.Body.Bytes(), &actual)
	is.Equal(actual, "success")
}
