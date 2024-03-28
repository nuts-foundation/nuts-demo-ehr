import OpenAPIClientAxios from "openapi-client-axios";
import definition from "./openapi-runtime.json";

function wrapClientCall(client, functionName, apiOptions, router) {
    return (parameters, payload, config) => {
        const sessionToken = localStorage.getItem("session")
        if (sessionToken) {
            if (!config) {
                config = {}
            }
            if (!config.headers) {
                config.headers = {}
            }
            config.headers['Authorization'] = `Bearer ${sessionToken}`
        }
        return client[functionName](parameters, payload, config)
            .then((result) => Promise.resolve(result))
            .catch((result) => {
                if (result.response.status === 401) {
                    // unauthorized
                    return router.push(apiOptions.forbiddenRoute)
                } else {
                    // If the response is a JSON error response, return the error text.
                    // Otherwise, return the more technical API client error message.
                    if (result.response.data && result.response.data.error) {
                        return Promise.reject(result.response.data.error)
                    } else {
                        return Promise.reject(result.message)
                    }
                }
            })
    }
}

export default {
    install: (app, apiOptions = {}) => {
        const client = new OpenAPIClientAxios({
            definition: definition,
        }).initSync();
        let proxy = {}
        for (const member in client) {
            if (typeof (client[member]) === 'function') {
                proxy[member] = wrapClientCall(client, member, apiOptions, app.config.globalProperties.$router)
            }
        }
        app.config.globalProperties.$api = proxy
    }
}
