_Please note that this project is currently a work in progress. Not everything
will fully work yet_

# Nuts demo EHR system

This application is pretending to be an electronic health record system. You can
use it to demo how healthcare professionals can work together by sharing
information with colleagues through the Nuts nodes.

It uses a FHIR server for the storage of patients, observataions and tasks.
We use the [Smart dev sandbox](https://github.com/smart-on-fhir/smart-dev-sandbox) with FHIR R4 and Synthea data
for a quick setup. You can generate a localized data set using the [Internaltion profiles for Synthea](https://github.com/synthetichealth/synthea-international).

This version is an updated version using vue.js as frontend framework and a Golang backend. It's based on the [nuts-registry-admin-demo](https://github.com/nuts-foundation/nuts-registry-admin-demo).

Go to the [master](https://github.com/nuts-foundation/nuts-registry-admin-demo/tree/master/) branch to find the previous version of the app.

**NOTE THAT THIS APPLICATION IS NOT INTENDED FOR USE WITH REAL MEDICAL
INFORMATION! IT IS IN NO WAY DEVELOPED TO BE SAFE, STABLE OR EVEN USABLE FOR
SUCH PURPOSE.**

## Building and running
### Production
To build for demo-production:

```shell
npm install
npm run build
go run .
```

This will serve the front end from the embedded filesystem.
### Development

During front-end development, you probably want to use the real filesystem and webpack in watch mode:

```shell
npm install
npm run watch
go run . live
```

The API and domain types are generated from the `api/api.yaml`.
```shell
npm run gen-api

oapi-codegen -generate server -package api api/api.yaml > api/generated.go
oapi-codegen -generate types -package domain -o domain/generated_types.go api/api.yaml
oapi-codegen -generate client,types -package auth -o client/auth/generated.go https://nuts-node.readthedocs.io/en/latest/_static/auth/v1.yaml
oapi-codegen -generate client,types -package vcr -o client/vcr/generated.go https://nuts-node.readthedocs.io/en/latest/_static/vcr/v1.yaml
oapi-codegen -generate client,types -package didman -o client/didman/generated.go -exclude-schemas OrganizationSearchResult https://nuts-node.readthedocs.io/en/latest/_static/didman/v1.yaml

```

### Docker
```shell
docker run -p 1304:1304 nutsfoundation/nuts-demo-ehr:main
```

#### Configuration
When running in Docker without a config file mounted at `/app/server.config.yaml` it will use the default configuration.

### Starting the _Smart On FHIR_ backend

The simplest way of starting up an out of the box FHIR backend is using the smart-on-fhir hapi server by running the following docker command:

```shell
docker run -p 4003:8080 smartonfhir/hapi-5:r3-synthea
```

This gives you a FHIR STU3 server with a generated patient population.
For other versions see the Smart on Health [docker repository](https://hub.docker.com/r/smartonfhir/hapi-5/tags?page=1&ordering=last_updated)

This server is part of a larger development kit which can be found [here](https://github.com/smart-on-fhir/smart-dev-sandbox#start-the-dev-sandbox).

#### FHIR configuration
The address of the FHIR endpoint can be configured in the `server.config.yaml` file.
The default value is set to:
```yaml
fhirserveraddr: "http://localhost:4003/hapi-fhir-jpaserver/fhir" 
```

### Nuts-node

The Demo-EHR needs a connection to a running Nuts node. The `customers.json` file also needs to be in sync with the DIDs known to the Nuts node.
You can use the [nuts-registry-admin-demo](https://github.com/nuts-foundation/nuts-registry-admin-demo) for setting up `customers.json`.

It's important to configure the Nuts node address in the `server.config.yaml`. The `nutsnodeaddr` must be used for this:

```yaml
nutsnodeaddr: "http://localhost:1323"
```

When using IRMA for authentication, the Nuts node will generate a QR code with an URL in it. This URL must be publicly accessible.
It can be configured IN THE NUTS NODE configuration file:

```yaml
auth:
  publicurl: http://5d6670ee3d46.eu.ngrok.io
```

The above example uses [ngrok](https://ngrok.io) to proxy a ngrok URL to localhost:1323.

## Technology Stack

Frontend framework is vue.js 3.x

Icons are from https://heroicons.com

CSS framework is https://tailwindcss.com
