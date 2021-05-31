export default {
  install: (app, apiOptions = {}) => {

    let {defaultOptions} = apiOptions

    const authHeader = () => {
      const sessionToken = localStorage.getItem("session")
      if (sessionToken) {
        return {'Authorization': `Bearer ${sessionToken}`}
      }
      return {}
    }
    let api = {}

    let httpMethods = ['get', 'post', 'put', 'delete']
    httpMethods.forEach((method) => {
      api[method] = (url, data = null, requestOptions = {}) => {
        const options = {
          ...defaultOptions,
          method: method.toUpperCase(),
          headers: {
            'Content-Type': 'application/json',
            ...authHeader()
          },
          ...requestOptions,
        }
        if (data) {
          options.body = JSON.stringify(data)
        }

        return fetch(url, options)
          .then((response) => {
            return response.json()
              .then((json) => {
                if (response.ok) {
                  return Promise.resolve(json)
                } else {
                  if (apiOptions.forbiddenRoute && response.status === 401) {
                    return app.config.globalProperties.$router.push(apiOptions.forbiddenRoute)
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
    })

    app.config.globalProperties.$api = api
  }
}