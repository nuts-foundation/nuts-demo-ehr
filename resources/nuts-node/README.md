# Nuts Node resources

These resources expose the Nuts Node API's defined [in the docs](https://nuts-documentation.readthedocs.io/en/latest/pages/API/index.html).

Each specific API has an OpenAPI specification. These resources each wrap an
OpenAPI client, generated based on those specifications. Note that we point to
the "live" OpenAPI specification files on Github here, which you wouldn't want
to do in production.

These resources also provide the application with an extra layer of abstraction,
so implementation details like "is the query sent in the path or in the POST
body" are abstracted away and we can add extra's like sorting the results.
