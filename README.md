
# walhallapiadaptor

`walhallapiadaptor` is intended to act as a fancy proxy for Walhall Core. It exposes the new style API for modules

It passes through the user's JWT without verifying it. (Verification is done implicitly based on whether or not
`WALHALL_API` allows for requests to be made with the supplied token.)

## Configuration
It takes the following environment variables:

| Variable | Description |
|---|---|
| `WALHALL_API_PREFIX` | The DNS name of the Walhall core API. (e.g. `http://api.walhall.io`) |
| `WALHALL_REGISTRY` | The DNS name of the default registry for Walhall. (Should be `registry.walhall.io`) |
| `PORT` | The port number the server should be exposed on. It defaults to `8080`. |

## Supported endpoints

| Method | Path Template | Description |
| --- | --- | ---|
| `GET` | `/orgs` | Returns a list of orgs a user is a member of |
| `GET` | `/orgs/{orgName}/modules` | Returns a list of modules in that organization |
| `POST` | `/orgs/{orgName}/modules/refresh` | *Temporary method* Initiates a sync of the modules for that org. |
| `GET` | `/orgs/{orgName}/modules/refresh` | *Temporary method* Gets the status of a sync for modules in an org. |

### Example response from GET /orgs/my-org/modules
    [
      {
        "name": "module-one",
        "source": "Github"
        "builds": [
          {
            "branch": "UNKNOWN",
            "commit": "UNKNOWN",
            "image": "registry.walhall.io/my-org/module-one:VERSION_ONE",
            "tags": [
              "VERSION_ONE"
            ]
          },
          {
            "branch": "UNKNOWN",
            "commit": "UNKNOWN",
            "image": "registry.walhall.io/my-org/module-one:VERSION_TWO",
            "tags": [
              "VERSION_TWO"
            ]
          }
        ]
      },
      {
        "name": "module-two",
        "source": "Github"
        "builds": [
          {
            "branch": "UNKNOWN",
            "commit": "UNKNOWN",
            "image": "registry.walhall.io/my-org/module-two:VERSION_ONE",
            "tags": [
              "VERSION_ONE"
            ]
          },
          {
            "branch": "UNKNOWN",
            "commit": "UNKNOWN",
            "image": "registry.walhall.io/my-org/module-two:VERSION_TWO",
            "tags": [
              "VERSION_TWO"
            ]
          }
        ]
      }
    ]


## Running locally

The service can be built with:

    $ go build humanitec.io/walhallapiadaptor/cmd/walhallapiadaptor

Tests can be run with:

    $ go test humanitec.io/walhallapiadaptor/cmd/walhallapiadaptor \
	    humanitec.io/walhallapiadaptor/internal/walhallapi

Mocks for the `humanitec.io/walhallapiadaptor/cmd/walhallapiadaptor` tests can be regenerated with:

    $ mockgen -source=../../internal/walhallapi/types.go -destination=walhallapier_mock.go -package=main WalhallAPIer
