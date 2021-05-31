_Please note that this project is currently a work in progress. Not everything
will fully work yet_

# Nuts demo EHR system

This application is pretending to be an electronic health record system. You can
use it to demo how healthcare professionals can work together by sharing
information with colleagues through the Nuts nodes.

This version is an updated version using vue.js as frontend framework and a Golang backend. It's based on the [nuts-registry-admin-demo](https://github.com/nuts-foundation/nuts-registry-admin-demo).

Go to the [master](https://github.com/nuts-foundation/nuts-registry-admin-demo/tree/master/) branch to find the previous version of the app.

**NOTE THAT THIS APPLICATION IS NOT INTENDED FOR USE WITH REAL MEDICAL
INFORMATION! IT IS IN NO WAY DEVELOPED TO BE SAFE, STABLE OR EVEN USABLE FOR
SUCH PURPOSE.**

## Learning from this application

If you're curious as to how this application interfaces with the Nuts node,
please take a look at [`resources/nuts-node`](resources/nuts-node), where we
define the different services and API calls that the Nuts node exposes. For
examples on how we then use those services, you can check out the [client APIs](client-api)
that the browser talks to to get things done. Mainly [`consent.js`](client-api/consent.js)
and [`organisation.js`](client-api/organisation.js). Also, we register our
applications on the Nuts node in the root [`index.js`](index.js).

