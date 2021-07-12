import Implementation from './api-impl'

function wrapClientCall(client, functionName, apiOptions, router) {
    return (parameters) => {
        return client[functionName](parameters).then((response) => {
            return response.json()
                .then((json) => {
                    if (response.ok) {
                        return Promise.resolve(json)
                    } else {
                        if (apiOptions.forbiddenRoute && response.status === 401) {
                            return router.push(apiOptions.forbiddenRoute)
                        } else {
                            return Promise.reject(json.error)
                        }
                    }
                }).catch(reason => {
                    // Handle 404 since it does not have content and the response.json() will fail.
                    if (response.status === 404) {
                        return Promise.reject(response.statusText)
                    }
                    // Handle 204 since it does not have content and the response.json() will fail.
                    if (response.status === 204) {
                        return Promise.resolve(response)
                    }
                    // Handle 201 since it might not have content and the response.json() will fail.
                    if (response.status === 201) {
                        return Promise.resolve(response)
                    }
                    return Promise.reject(reason)
                })
        })
    }
}

export default {
    install: (app, apiOptions = {}) => {
        const client = new Implementation({
            cors: true,
            endpoint: '.',
            securityHandlers: {
                "bearerAuth": (headers, params) => {
                    const sessionToken = localStorage.getItem("session")
                    headers["Authorization"] = `Bearer ${sessionToken}`
                    return true
                }
            }
        });

        // Wrap all generated OpenAPI client functions with code for handling error and JSON responses
        let proxy = {}
        for (const member in client) {
            if (typeof (client[member]) === 'function') {
                proxy[member] = wrapClientCall(client, member, apiOptions, app.config.globalProperties.$router)
            }
        }

        app.config.globalProperties.api = proxy
    }
}
