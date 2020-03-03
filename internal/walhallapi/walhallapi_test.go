package walhallapi

import (
	"net/http"
	"testing"

	"github.com/matryer/is"
	"humanitec.io/walhallapiadaptor/internal/testutil"
)

const exampleJWT = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJpc3MiOiJ3YWxoYWxsIiwiZXhwIjoxNTgwMjAxNDc1LCJpYXQiOjE1ODAxMTUwNzUsInVzZXJfdXVpZCI6IjBiNjE4NTc5LWY1NDYtNDMzOC05ZWNlLWExYzk4MWY5MGM4MCIsIm9yZ2FuaXphdGlvbl91dWlkcyI6WyJhNzlkOWU5OS00NzZkLTRlMjktYTVmMi02MDEwMmE1ZmZmMjkiLCJmMzNmMDEzZS1lNTMyLTRiMjctOTU4ZS01MDIyMGExOGEyYmQiXSwidXNlcm5hbWUiOiJjaHJpc2h1bWFuaXRlYyIsInNjb3BlIjoicmVhZCB3cml0ZSJ9.aGHKK1G5POgKtGaUB8duGfnt_F_RXeTsnzx3gOLOf3w3YqVgB440IIFUreFYfxWpJtyOzBja1ytRLPHOFFis3kwBqOlWsfDqrffwzaur-ESBAC4DE_rEtw52UXR4aEozu88bfuk3CpMQDy66ZVNbPfrUMPh9Mh51a_1E74DZt8Tm3O66afc9XT9Xkv_RwFncPWuLIGzxAj-tyIFDGfYxqIPy61NySrOsvZ88G4x6yU5w4fS1gMRkwjKCepId72YdcScTEp5Qr4O6yHHRDTd43zJ06CmDniy1FEhIgBdD6zyJZhvXn5dr_aJSUMlduJSzVu32A1owm54Qg77JKhbuXA"

const getUserResponse = `{
  "user_uuid": "0b618579-f546-4338-9ece-a1c981f90c80",
  "organizations": [
    {
      "id": 84,
      "organization_uuid": "a79d9e99-476d-4e29-a5f2-60102a5fff29",
      "name": "chrishumanitec",
      "description": null,
      "organization_url": null,
      "company_name": null,
      "team_size": null,
      "create_date": "2019-07-15T18:09:17.368713+02:00",
      "edit_date": "2019-07-15T18:09:17.368722+02:00",
      "subscription_id": null,
      "used_seats": 0,
      "profile_name": null,
      "industry": []
    },
    {
      "id": 85,
      "organization_uuid": "f33f013e-e532-4b27-958e-50220a18a2bd",
      "name": "corporate-org",
      "description": null,
      "organization_url": null,
      "company_name": null,
      "team_size": null,
      "create_date": "2019-07-15T18:09:22.117490+02:00",
      "edit_date": "2019-07-15T18:09:22.117500+02:00",
      "subscription_id": null,
      "used_seats": 0,
      "profile_name": null,
      "industry": []
    }
  ]
}`

const getListAppsResponse = `{
  "count": 2,
  "next": null,
  "previous": null,
  "results": [
    {
      "app_uuid": "10a1604d-da69-4e12-a5c6-ac5fad87ae62",
      "id": "10a1604d-da69-4e12-a5c6-ac5fad87ae62",
      "environments": [
        "fa9852ef-963c-45a8-a420-0f099543c989"
      ],
      "app_url": null,
      "name": "test-app-one",
      "description": null,
      "domain_name": null,
      "github_repo": null,
      "create_date": "2020-01-27T10:20:42.590626+01:00",
      "edit_date": "2020-01-27T10:20:42.590634+01:00",
      "logo": null
    },
    {
      "app_uuid": "c0859864-3f2c-40c9-bff9-5a227a31d379",
      "id": "c0859864-3f2c-40c9-bff9-5a227a31d379",
      "environments": [
        "acfe54c3-60ae-41e5-a8ea-d34951aa4b7f",
        "fdb57593-e794-4786-9b38-c59d44cf5d6c"
      ],
      "app_url": null,
      "name": "test-app-two",
      "description": null,
      "domain_name": null,
      "github_repo": null,
      "create_date": "2020-01-21T15:49:37.102272+01:00",
      "edit_date": "2020-01-21T15:49:37.102279+01:00",
      "logo": null
    }
  ]
}`

const getListModules = `{
  "count": 2,
  "next": null,
  "previous": null,
  "results": [
    {
      "id": "d4705bb7-abc5-4efd-b7e3-bbd041d4aebd",
      "module_uuid": "d4705bb7-abc5-4efd-b7e3-bbd041d4aebd",
      "name": "eve-demo",
      "description": "Product Inventory Frontend. Expects to talk to a zed-demo module in-cluster.",
      "endpoint": null,
      "framework_id": null,
      "is_system_module": false,
      "github_repo": "corporate-org/eve-demo",
      "image": "eve-demo",
      "service_type": "backend",
      "status": "internal",
      "create_date": "2019-10-24T11:08:15.381517+02:00",
      "edit_date": "2020-01-21T15:48:49.663251+01:00",
      "upstream": null,
      "versions": [
        {
          "id": 11925,
          "version_uuid": "ddaaed43-3d16-4c91-ad3a-ca62ed2cf911",
          "version": "0.0.3",
          "repository_tag": "0.0.3",
          "core_version": null,
          "status": "published",
          "create_date": "2019-10-24T16:16:11+02:00",
          "edit_date": "2020-01-03T13:37:50.351486+01:00"
        },
        {
          "id": 11776,
          "version_uuid": "dcf9d403-6bb0-42fb-8bf5-2fa7ea0d7926",
          "version": "0.0.2",
          "repository_tag": "0.0.2",
          "core_version": null,
          "status": "published",
          "create_date": "2019-10-24T13:49:27+02:00",
          "edit_date": "2020-01-03T13:37:50.804946+01:00"
        },
        {
          "id": 11738,
          "version_uuid": "47218098-98e6-4c12-b91b-71a8f23dc1a3",
          "version": "0.0.1",
          "repository_tag": "0.0.1",
          "core_version": null,
          "status": "published",
          "create_date": "2019-10-24T11:04:53+02:00",
          "edit_date": "2020-01-03T13:37:51.222580+01:00"
        }
      ],
      "documentation": "https://github.com/corporate-org/eve-demo",
      "organization_id": 85,
      "is_deployed": true,
      "repository_status": "ready"
    },
    {
      "id": "f51d73c6-5a15-4efa-aa4c-21822437ef94",
      "module_uuid": "f51d73c6-5a15-4efa-aa4c-21822437ef94",
      "name": "product-be",
      "description": "Backend for Product Listing app",
      "endpoint": null,
      "framework_id": null,
      "is_system_module": false,
      "github_repo": "corporate-org/product-be",
      "image": "product-be",
      "service_type": "backend",
      "status": "internal",
      "create_date": "2019-12-16T13:27:37.064844+01:00",
      "edit_date": "2020-01-21T15:48:57.416658+01:00",
      "upstream": null,
      "versions": [
        {
          "id": 20129,
          "version_uuid": "80d10e16-34b2-4757-b794-562d4ae0770f",
          "version": "v1.0",
          "repository_tag": "v1.0",
          "core_version": null,
          "status": "published",
          "create_date": "2019-12-18T12:10:31+01:00",
          "edit_date": "2020-01-21T15:48:58.532027+01:00"
        }
      ],
      "documentation": "https://github.com/corporate-org/product-be",
      "organization_id": 85,
      "is_deployed": false,
      "repository_status": "ready"
    }
  ]
}`

const getEnvrionmentFromApp = `{
  "count": 1,
  "next": null,
  "previous": null,
  "results": [
    {
      "env_uuid": "fa9852ef-963c-45a8-a420-0f099543c989",
      "id": "fa9852ef-963c-45a8-a420-0f099543c989",
      "logic_module_versions": [
        {
          "id": 18862,
          "logic_module": {
            "id": "44bc53dd-142a-41a4-9d29-896f5fb3f0d0",
            "module_uuid": "44bc53dd-142a-41a4-9d29-896f5fb3f0d0",
            "name": "demo-be",
            "description": "Backend of a Product Inventory Service",
            "endpoint": null,
            "framework_id": null,
            "is_system_module": false,
            "github_repo": "corporate-org/demo-be",
            "image": "demo-be",
            "service_type": "backend",
            "status": "internal",
            "create_date": "2020-01-03T13:37:59.449119+01:00",
            "edit_date": "2020-01-27T23:14:47.227665+01:00",
            "upstream": null,
            "versions": [
              {
                "id": 18862,
                "version_uuid": "d48ded59-4b56-4b1b-94b7-9757856952a4",
                "version": "1.0",
                "repository_tag": "1.0",
                "core_version": null,
                "status": "published",
                "create_date": "2020-01-03T13:34:09+01:00",
                "edit_date": "2020-01-03T13:38:00.827043+01:00"
              }
            ],
            "documentation": "https://github.com/corporate-org/demo-be",
            "organization_id": 85,
            "is_deployed": true,
            "repository_status": "ready"
          },
          "version_uuid": "d48ded59-4b56-4b1b-94b7-9757856952a4",
          "version": "1.0",
          "repository_tag": "1.0",
          "core_version": null,
          "status": "published",
          "create_date": "2020-01-03T13:34:09+01:00",
          "edit_date": "2020-01-03T13:38:00.827043+01:00"
        }
      ],
      "logic_module_version_ingress": {
        "d48ded59-4b56-4b1b-94b7-9757856952a4": "development-romerohuffnichols.newapp.io"
      },
      "cluster": {
        "id": 9,
        "k8s_cluster_uuid": "ad7a227f-4516-41d7-9cb4-c16e8a5e661d",
        "name": "cluster-2",
        "is_default": true
      },
      "name": "Development",
      "is_production": false,
      "create_date": "2020-01-27T23:12:38.699450+01:00",
      "edit_date": "2020-01-27T23:19:27.837229+01:00",
      "deployment_status": "success"
    }
  ]
}`

const getConfigsPerModulePerEnv = `{
  "count": 4,
  "next": null,
  "previous": null,
  "results": [
    {
      "id": 26013,
      "configuration_uuid": "58d7ed31-cefb-41d5-9c52-956111361651",
      "name": "demobe-config-map",
      "specification": {
        "data": {
          "EXAMPLE_VAR": "example",
          "OTHER_EXAMPLE": "other"
        }
      },
      "type": "config_map",
      "status": "done",
      "environment": "fa9852ef-963c-45a8-a420-0f099543c989",
      "logic_module_version": 18862,
      "create_date": "2020-01-27T23:12:38.721794+01:00",
      "edit_date": "2020-01-27T23:19:16.439895+01:00"
    },
    {
      "id": 26014,
      "configuration_uuid": "e2bb3a0f-beea-434c-99b7-a0f98c155793",
      "name": "demobe-container",
      "specification": {
        "name": "demobe-container",
        "image": "registry.walhall.io/corporate-org/demo-be:1.0",
        "ports": [
          {
            "container_port": 8080
          }
        ],
        "resources": {
          "limits": {
            "cpu": "1000m",
            "memory": "1024Mi"
          },
          "requests": {
            "cpu": "250m",
            "memory": "256Mi"
          }
        }
      },
      "type": "container",
      "status": "done",
      "environment": "fa9852ef-963c-45a8-a420-0f099543c989",
      "logic_module_version": 18862,
      "create_date": "2020-01-27T23:12:38.733179+01:00",
      "edit_date": "2020-01-27T23:19:16.794839+01:00"
    },
    {
      "id": 26017,
      "configuration_uuid": "d8b6eb44-8533-427d-b0fe-d881950fe669",
      "name": "demobe-ingress",
      "specification": {
        "spec": {
          "tls": [
            {
              "hosts": [
                "development-romerohuffnichols.newapp.io"
              ],
              "secret_name": "wildcard-newapp-io-tls"
            }
          ],
          "rules": [
            {
              "host": "development-romerohuffnichols.newapp.io",
              "http": {
                "paths": [
                  {
                    "backend": {
                      "service_name": "demobe-service",
                      "service_port": 80
                    }
                  }
                ]
              }
            }
          ]
        },
        "metadata": {
          "name": "demobe-ingress",
          "namespace": "fa9852ef-963c-45a8-a420-0f099543c989",
          "annotations": {
            "kubernetes.io/ingress.class": "nginx",
            "nginx.org/client-max-body-size": "20m",
            "ingress.kubernetes.io/ssl-redirect": "true",
            "nginx.ingress.kubernetes.io/enable-cors": "true",
            "nginx.ingress.kubernetes.io/cors-allow-origin": "https://app.walhall.io,newapp.io"
          }
        }
      },
      "type": "ingress",
      "status": "done",
      "environment": "fa9852ef-963c-45a8-a420-0f099543c989",
      "logic_module_version": 18862,
      "create_date": "2020-01-27T23:14:43.701323+01:00",
      "edit_date": "2020-01-27T23:14:58.630281+01:00"
    },
    {
      "id": 26015,
      "configuration_uuid": "fc7eb08a-e9d5-42a2-ba3b-cfc438b0163a",
      "name": "demobe-service",
      "specification": {
        "spec": {
          "type": "ClusterIP",
          "ports": [
            {
              "name": "http",
              "port": 80,
              "protocol": "TCP",
              "target_port": 8080
            }
          ],
          "selector": {
            "app": "demobe"
          }
        },
        "status": {
          "load_balancer": {}
        },
        "metadata": {
          "name": "demobe-service",
          "namespace": "fa9852ef-963c-45a8-a420-0f099543c989"
        }
      },
      "type": "service",
      "status": "done",
      "environment": "fa9852ef-963c-45a8-a420-0f099543c989",
      "logic_module_version": 18862,
      "create_date": "2020-01-27T23:12:38.784714+01:00",
      "edit_date": "2020-01-27T23:14:47.450426+01:00"
    }
  ]
}`

func TestListOrgs(t *testing.T) {
	is := is.New(t)
	client := testutil.NewFakeDoer(t)
	client.HandleRequest("GET", "/api/walhalluser/0b618579-f546-4338-9ece-a1c981f90c80", http.StatusOK, []byte(getUserResponse), t)

	helper, err := New("http://api.walhall.io", exampleJWT, client)
	if err != nil {
		t.Error(err)
	}
	orgs, err := helper.ListOrgs()
	is.NoErr(err)

	is.Equal("f33f013e-e532-4b27-958e-50220a18a2bd", orgs["corporate-org"])

}

func TestListApps(t *testing.T) {
	is := is.New(t)
	client := testutil.NewFakeDoer(t)
	client.HandleRequest("GET", "/api/walhalluser/0b618579-f546-4338-9ece-a1c981f90c80", http.StatusOK, []byte(getUserResponse), t)
	client.HandleRequest("GET", "/api/application?limit=100&organization_uuid=f33f013e-e532-4b27-958e-50220a18a2bd", http.StatusOK, []byte(getListAppsResponse), t)

	helper, err := New("http://api.walhall.io", exampleJWT, client)
	if err != nil {
		t.Error(err)
	}
	apps, err := helper.ListApps("corporate-org")
	is.NoErr(err)

	is.Equal("10a1604d-da69-4e12-a5c6-ac5fad87ae62", apps["test-app-one"])
	is.Equal("c0859864-3f2c-40c9-bff9-5a227a31d379", apps["test-app-two"])

}

func TestListModules(t *testing.T) {
	is := is.New(t)
	client := testutil.NewFakeDoer(t)
	client.HandleRequest("GET", "/api/walhalluser/0b618579-f546-4338-9ece-a1c981f90c80", http.StatusOK, []byte(getUserResponse), t)
	client.HandleRequest("GET", "/api/logicmodule?organization=f33f013e-e532-4b27-958e-50220a18a2bd&limit=50&status=internal", http.StatusOK, []byte(getListModules), t)

	helper, err := New("http://api.walhall.io", exampleJWT, client)
	if err != nil {
		t.Error(err)
	}
	modules, err := helper.ListModules("corporate-org")
	is.NoErr(err)

	is.Equal(2, len(modules)) // Expecting 2 modules
	is.Equal("eve-demo", modules[0].Name)
	is.Equal(3, len(modules[0].Versions))

}

/*
func TestGetEnvironmentAsDeploymentSet(t *testing.T) {
	is := is.New(t)
	client := NewFakeDoer(t)
	client.HandleRequest("GET", "/api/walhalluser/0b618579-f546-4338-9ece-a1c981f90c80", http.StatusOK, []byte(getUserResponse), t)
	client.HandleRequest("GET", "/api/application?limit=100&organization_uuid=f33f013e-e532-4b27-958e-50220a18a2bd", http.StatusOK, []byte(getListAppsResponse), t)
	client.HandleRequest("GET", "/api/environments?application=10a1604d-da69-4e12-a5c6-ac5fad87ae62", http.StatusOK, []byte(getEnvrionmentFromApp), t)
	client.HandleRequest("GET", "/api/configuration?logic_module_version=18862&environment=fa9852ef-963c-45a8-a420-0f099543c989", http.StatusOK, []byte(getConfigsPerModulePerEnv), t)
	client.HandleBodyRequest("POST", "/orgs/corporate-org/apps/test-app-one/sets/0", http.StatusOK, []byte(`"0123456789abcdef09765"`), func(body io.ReadCloser) bool {
		var delta depset.Delta
		decoder := json.NewDecoder(body)
		decoder.Decode(&delta)
		return reflect.DeepEqual(delta, depset.Delta{
			Modules: depset.ModuleDeltas{
				Add: map[string]depset.ModuleSpec{
					"demo-be": depset.ModuleSpec{
						"config-map": map[string]interface{}{
							"EXAMPLE_VAR":   "example",
							"OTHER_EXAMPLE": "other",
						},
						"ingress": true,
					},
				},
			},
		})
	}, t)
	helper, err := New("http://api.walhall.io", exampleJWT, client)
	if err != nil {
		t.Error(err)
	}
	setID, err := helper.GetEnvironmentAsDeploymentSet("corporate-org", "test-app-one", "Development")
	is.NoErr(err)
	is.Equal("0123456789abcdef09765", setID)
}*/
