// # check access via token introspection as described by https://www.nginx.com/blog/validating-oauth-2-0-access-tokens-nginx/
function introspectAccessToken(r) {
    // check if the Authorization header is present and long enough
    if (!r.headersIn['Authorization'] || r.headersIn['Authorization'].length < 7) {
        r.return(403);
        return
    }
    // strip the first 5 or 7 chars
    const token = "token=" + r.headersIn['Authorization'].substring(7);
    // make a subrequest to the introspection endpoint
    r.subrequest("/_oauth2_introspect",
        { method: "POST", body: token},
        function(reply) {
            if (reply.status === 200) {
                const introspection = JSON.parse(reply.responseBody);
                if (introspection.active === true) {
                    //dpop(r, introspection.cnf)
                    r.headersOut['X-Userinfo'] = btoa(reply.responseBody);
                    r.return(200);
                } else {
                    r.return(403);
                }
            } else {
                r.return(500);
            }
        }
    );
}

// call the PDP to check the access
function authorize(r) {
    // const xUserinfo = r.headersIn['X-Userinfo'];
    // const requestLine = r.request
    const input =
        JSON.stringify({
            "input": {
                "request": {
                    "method": r.variables.request_method,
                    "path": r.variables.request_uri, // original non-normalized request_uri, may need some processing in more complex situations
                    "headers": {
                        "X-Userinfo": r.headersIn["X-Userinfo"]
                    }
                }
            }
        });
    r.subrequest("/_oauth2_authorize",
        { method: "POST", body: input},
        function(reply) {
            if (reply.status === 200) {
                r.error(reply.responseBody);
                const authResult = JSON.parse(reply.responseBody);
                if (authResult.result.allow === true) {
                    r.return(200);
                } else {
                    r.return(403);
                }
            } else {
                r.error(reply.responseBody);
                r.return(500);
            }
        }
    );
}

export default { introspectAccessToken, authorize };