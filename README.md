_Please note that this project is currently a work in progress. Not everything
will fully work yet_

# Nuts demo EHR system

This application is pretending to be an electronic health record system. You can
use it to demo how healthcare professionals can work together by sharing
information with colleagues through the Nuts nodes.

Also, you can use it as "the other party" in implementing Nuts in your own EHR
systems. We gave a live demo of this application in one of our meetings, which
[you can see here](https://www.youtube.com/watch?v=TONCu0AHPWs) (in Dutch).

**NOTE THAT THIS APPLICATION IS NOT INTENDED FOR USE WITH REAL MEDICAL
INFORMATION! IT IS IN NO WAY DEVELOPED TO BE SAFE, STABLE OR EVEN USABLE FOR
SUCH PURPOSE.**

## Quick start guide

```bash
git clone git@github.com:nuts-foundation/nuts-demo-ehr.git
cd nuts-demo-ehr
npm install
npm start
```

Note that it needs a running redis server for session persistence. You can start one by running a docker node:
```bash
docker run -p 127.0.0.1:6379:6379 --rm --name demo-ehr-session-redis redis
```

You should now have three instances of this EHR running on:

* http://localhost:8000 ⸺ Verpleeghuis de Nootjes
* http://localhost:8001 ⸺ Huisartsenpraktijk Nootenboom
* http://localhost:8002 ⸺ Medisch Centrum Noot aan de Man

Also, as a bonus, you can display two or all three side by side by going to:

* http://localhost:8000/duo.html ⸺ Shows the applications on ports 8000 and 8001
* http://localhost:8000/triple.html ⸺ Shows all three applications

### Configuring the application(s)

You can find the configuration files for all three applications in the `config`
directory. You may need to edit these files to point to the right Nuts node(s).
If you followed the [Setup a local Nuts network](https://nuts-documentation.readthedocs.io/en/latest/pages/getting_started/local_network.html#setup-a-local-nuts-network)
instructions, http://localhost:8000 (Verpleeghuis de Nootjes) will connect to Nuts
node 'Bundy' and the other two will connect to node 'Dahmer'.

You can also change port numbers, organisation details and default health
records in the config files.

### Configuring the Nuts node

### IRMA

You'll need to change a few things in the Nuts node config if you want to use
the IRMA flows. In the `nuts-network-local/config/bundy/nuts.yaml` and
`nuts-network-local/config/dahmer/nuts.yaml` files you will find a line that
reads:

```yaml
auth:
  publicUrl: https://example.org
```

You will need to change this to a URL that both your browser and your phone can
connect to the Nuts node on. So this can be something like:

```yaml
auth:
  publicUrl: http://192.168.1.xx:11323
```

Or you can use a service like ngrok to proxy requests to your local machine.

You may also want to set this value to allow you to test using
[demo attributes](https://privacybydesign.foundation/attribute-index/en/irma-demo.html):

```yaml
auth:
  irmaSchemeManager: irma-demo
```

You will need to restart your Nuts nodes to enable these changes.

#### NATS events

In order to receive consent events from the Nuts node you will have to allow the
demo EHR to talk to NATS. You do this by adding these port forwards to
`nuts-network-local/docker-compose.yml`:

```yaml
bundy-nuts-service-space:
  ...
  ports:
    - "11323:1323"
    - "11324:4222" # <-- Add this
  ...
dahmer-nuts-service-space:
  ...
  ports:
    - "21323:1323"
    - "21324:4222" # <-- Add this
  ...
```

You will need to restart your Nuts nodes to enable these changes.

### Adding to the Nuts register

If you want to allow the applications to find each other and exchange data, you
will have to add them to the Nuts registry. Again, if you followed
[Setup a local Nuts network](https://nuts-documentation.readthedocs.io/en/latest/pages/getting_started/local_network.html#setup-a-local-nuts-network),
you can find your registry in `nuts-network-local/config/registry`.

#### 1. Add the organisations

To make this process easier the applications will output their registry
information on startup. Add that information to your registry's
`organisations.json`.

#### 2. Add the endpoints

You can add the locations of the APIs to the `endpoints.json` file as endpoints
of the type `urn:oid:1.3.6.1.4.1.54851.2:demo-ehr`. Also, each Nuts node that
can receive consent needs an endpoint of the type `urn:nuts:endpoint:consent`.

So for each application add this endpoint to `endpoints.json`:

```json
{
  "endpointType": "urn:oid:1.3.6.1.4.1.54851.2:demo-ehr",
  "identifier": "0e906b06-db48-452c-bb61-559f239a06ca",
  "status": "active",
  "version": "0.1.0",
  "URL": "http://localhost:8000/external/patient"
}
```

Make sure you give each one a unique identifier and have it point to the right
URL. Also, make sure both Nuts nodes have a consent endpoint (should be okay if
you're using `nuts-network-local`).

#### 3. Connect the endpoints to organisations

Connect your endpoints to organizations in the `endpoint_organizations.json`
file like this:

```json
{
  "status": "active",
  "organization": "urn:oid:2.16.840.1.113883.2.4.6.1:12345678",
  "endpoint": "0e906b06-db48-452c-bb61-559f239a06ca"
}
```

Make sure you add two entries for each organisation, one for the API and one for
the Nuts node consent endpoint.

Note that for hot reloading of the registry to be triggered, every one of these
three file needs to be touched.

## Learning from this application

If you're curious as to how this application interfaces with the Nuts node,
please take a look at [`resources/nuts-node`](resources/nuts-node), where we
define the different services and API calls that the Nuts node exposes. For
examples on how we then use those services, you can check out the [client APIs](client-api)
that the browser talks to to get things done. Mainly [`consent.js`](client-api/consent.js)
and [`organisation.js`](client-api/organisation.js). Also, we register our
applications on the Nuts node in the root [`index.js`](index.js).
