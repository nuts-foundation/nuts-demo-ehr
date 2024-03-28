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

        // TODO: Might have to handle the old functionality as well:
        // return client[functionName](parameters).then((response) => {
        //     return response.json()
        //         .then((json) => {
        //             if (response.ok) {
        //                 return Promise.resolve(json)
        //             } else {
        //                 if (apiOptions.forbiddenRoute && response.status === 401) {
        //                     return router.push(apiOptions.forbiddenRoute)
        //                 } else {
        //                     return Promise.reject(json.error)
        //                 }
        //             }
        //         }).catch(reason => {
        //             // Handle 404 since it does not have content and the response.json() will fail.
        //             if (response.status === 404) {
        //                 return Promise.reject(reason)
        //             }
        //             // Handle 204 since it does not have content and the response.json() will fail.
        //             if (response.status === 204) {
        //                 return Promise.resolve(response)
        //             }
        //             // Handle 201 since it might not have content and the response.json() will fail.
        //             if (response.status === 201) {
        //                 return Promise.resolve(response)
        //             }
        //             return Promise.reject(reason)
        //         })
        // })
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
