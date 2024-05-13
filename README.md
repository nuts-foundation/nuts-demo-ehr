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

Versions are in sync with the Nuts node version. The main branch uses the master branch of the Nuts node.
Older versions have a `vX` branch.

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
make gen-api
```

### Docker
```shell
docker run -p 1304:1304 nutsfoundation/nuts-demo-ehr:main
```

#### Configuration
When running in Docker without a config file mounted at `/app/server.config.yaml` it will use the default configuration.

#### TLS

To allow Demo EHR to query FHIR servers and eOverdracht notification endpoints which require a client certificate (required according to the Bolt),
you need to configure `tls.client.certificate` and `tls.client.key` to point to the respective files.
Use the same certificate you're using for your Nuts node.

There's no need to configure the truststore: Demo EHR skips verification of the server certificate (it's a demo application after all).

### Starting the HAPI FHIR server backend

The simplest way of starting up an out of the box FHIR backend is using the HAPI FHIR server by running the following docker command:

```shell
docker run -p 8080:8080 -e hapi.fhir.fhir_version=DSTU3 -e hapi.fhir.partitioning.allow_references_across_partitions=false hapiproject/hapi:v5.4.1
```

Configuration explanation:
- `hapi.fhir.fhir_version=DSTU3` indicates FHIR version STU3 is used
- `hapi.fhir.partitioning.allow_references_across_partitions=false` signals HAPI server to enable partitioning, which allows multi-tenancy.

### FHIR server type

If you're using the HAPI FHIR docker image or any other HAPI FHIR server with support for multi-tenancy you should set the `fhir.server.type` option to: `hapi-multi-tenant` otherwise choose either `hapi` (for a single-tenant HAPI FHIR server) or `other`.

### Nuts-node

The Demo-EHR needs a connection to a running Nuts node. The `customers.json` file also needs to be in sync with the DIDs known to the Nuts node.
You can use the [nuts-registry-admin-demo](https://github.com/nuts-foundation/nuts-registry-admin-demo) for setting up `customers.json`.

It's important to configure the Nuts node address in the `server.config.yaml`. The `nutsnodeaddr` must be used for this:

```yaml
nutsnodeaddr: "http://localhost:1323"
```

## Technology Stack

Frontend framework is vue.js 3.x

Icons are from https://heroicons.com

CSS framework is https://tailwindcss.com
