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
If you start the local development network using `make` in https://github.com/nuts-foundation/nuts-node/tree/master/development,
http://localhost:8000 (Verpleeghuis de Nootjes) will connect to `node1` and the other two will connect to `node2`.

You can also change port numbers, organisation details and default health
records in the config files.

### Configuring the Nuts node

### IRMA

You'll need to change a few things in the Nuts node config if you want to use
the IRMA flows. In the `nuts-node/development/nuts.yaml` you will need to specify the public URL that both your browser
and your phone can connect to the Nuts node on. So this can be something like:

```yaml
auth:
  publicurl: http://192.168.1.xx:11323
```

Or you can use a service like ngrok to proxy requests to your local machine.

You may also want to set the IRMA scheme manager to `irma-demo` to allow you to test using
[demo attributes](https://privacybydesign.foundation/attribute-index/en/irma-demo.html):

```yaml
auth:
  irma:
    schememanager: irma-demo
```

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

## Learning from this application

If you're curious as to how this application interfaces with the Nuts node,
please take a look at [`resources/nuts-node`](resources/nuts-node), where we
define the different services and API calls that the Nuts node exposes. For
examples on how we then use those services, you can check out the [client APIs](client-api)
that the browser talks to to get things done. Mainly [`consent.js`](client-api/consent.js)
and [`organisation.js`](client-api/organisation.js). Also, we register our
applications on the Nuts node in the root [`index.js`](index.js).
