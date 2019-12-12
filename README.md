_Please note that this project is currently a work in progress. Not everything
will fully work yet_

# Nuts demo EHR system

This application is pretending to be an electronic health record system. You can
use it to demo how healthcare professionals can work together by sharing
information with colleagues through the Nuts nodes.

Also, you can use it as "the other party" in implementing Nuts in your own EHR
systems.

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

You should now have three instances of this EHR running on:

* http://localhost:80 ⸺ Verpleeghuis de Nootjes
* http://localhost:81 ⸺ Huisartsenpraktijk Nootenboom
* http://localhost:82 ⸺ Medisch Centrum Noot aan de Man

Also, as a bonus, you can display two or all three side by side by going to:

* http://localhost/duo.html ⸺ Shows the applications on ports 80 and 82
* http://localhost/triple.html ⸺ Shows all three applications

### Configuring the application(s)

You can find the configuration files for all three applications in the `config`
directory. You may need to edit these files to point to the right Nuts node(s).
If you followed the [Setup a local Nuts network](https://nuts-documentation.readthedocs.io/en/latest/pages/getting_started/local_network.html#setup-a-local-nuts-network)
instructions, http://localhost:80 (Verpleeghuis de Nootjes) will connect to Nuts
node 'Bundy' and the other two will connect to node 'Dahmer'.

You can also change port numbers, organisation details and default health
records in the config files.

### Adding to the Nuts register

If you want to allow the applications to find each other and exchange data, you
will have to add them to the Nuts registry. Again, if you followed
[Setup a local Nuts network](https://nuts-documentation.readthedocs.io/en/latest/pages/getting_started/local_network.html#setup-a-local-nuts-network),
you can find your registry in `nuts-network-local/config/registry`.

To make this process easier the applications will output their registry
information on startup. Add that information to your registry's
`organisations.json`.

## Learning from this application

If you're curious as to how this application interfaces with the Nuts node,
please take a look at [`resources/nuts-node`](resources/nuts-node), where we
define the different services and API calls that the Nuts node exposes. For
examples on how we then use those services, you can check out the [client APIs](client-api)
that the browser talks to to get things done. Mainly [`consent.js`](client-api/consent.js)
and [`organisation.js`](client-api/organisation.js). Also, we register our
applications on the Nuts node in the root [`index.js`](index.js).
