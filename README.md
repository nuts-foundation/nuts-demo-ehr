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

## WIP: complete docker compose setup with 2 instances

# Nodes need to find each other by external URL because they need to resolve the web:did
Solution: handout 172.90.0.2 to loadbalancer
Set default range to:
networks:
    default:
        ipam:
            config:
                - subnet: 172.90.0.0/16
### After clone

- execute `./generate.sh` in `docker-compose/lb/tls/`
- load `docker-compose/lb/tls/ca.pem` into keychain/local certs
- add `left.local`, `node.left.local`, `admin.left.local`, `right.local`, `admin.right.local`  and `node.right.local` to `/etc/hosts` (127.0.0.1)
- execute `make docker`
- execute `./setup.sh` (this will also start all containers in the docker-compose setup)

### setup.sh
the setup script executes the following steps:
- use https://admin.left.local and add did:web:left.local:iam:left
- issue an NutsOrganizationCredential for this DID from this DID
- use https://admin.right.local and add did:web:right.local:iam:right
- issue an NutsOrganizationCredential for this DID from this DID
- enable services for discovery
- add services to the DID document using curl statement from below
- wait

curl statement to add services to DID document:
```shell
docker exec nuts-demo-ehr-node-left-1 curl -X POST "http://localhost:8081/internal/vdr/v2/did/did:web:node.left.local:iam:left/service" -H  "Content-Type: application/json" -d '{"type": "eOverdracht-sender","serviceEndpoint": {"auth": "https://node.left.local/oauth2/did:web:node.left.local:iam:left/authorize","fhir": "https://left.local/fhir/1"}}'
docker exec nuts-demo-ehr-node-left-1 curl -X POST "http://localhost:8081/internal/vdr/v2/did/did:web:node.left.local:iam:left/service" -H  "Content-Type: application/json" -d '{"type": "eOverdracht-receiver","serviceEndpoint": {"auth": "https://node.left.local/oauth2/did:web:node.left.local:iam:left/authorize","notification": "https://left.local/web/external/transfer/notify"}}'
docker exec nuts-demo-ehr-node-right-1 curl -X POST "http://localhost:8081/internal/vdr/v2/did/did:web:node.right.local:iam:right/service" -H  "Content-Type: application/json" -d '{"type": "eOverdracht-sender","serviceEndpoint": {"auth": "https://node.right.local/oauth2/did:web:node.right.local:iam:right/authorize","fhir": "https://right.local/fhir/1"}}'
docker exec nuts-demo-ehr-node-right-1 curl -X POST "http://localhost:8081/internal/vdr/v2/did/did:web:node.right.local:iam:right/service" -H  "Content-Type: application/json" -d '{"type": "eOverdracht-receiver","serviceEndpoint": {"auth": "https://node.right.local/oauth2/did:web:node.right.local:iam:right/authorize","notification": "https://right.local/web/external/transfer/notify"}}'
```

### Run
- docker compose up
- goto: https://left.local and https://right.local