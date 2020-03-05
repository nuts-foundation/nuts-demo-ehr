const OpenAPIClientAxios = require('openapi-client-axios').default

module.exports = {
  call: ({ baseURL, definition }) => {
    const auth = new OpenAPIClientAxios({
      definition,
      strict: false,
      validate: false,
      axiosConfigDefaults: { baseURL }
    })
    auth.init()

    return async function (method, params = null, body = null, headers = null) {
      const config = {
        headers: {
          Accept: 'application/json',
          ...headers
        }
      }
      try {
        const client = await auth.getClient()
        const result = await client[method](params, body, config)
        return result.data
      } catch (e) {
        throw (e)
      }
    }
  },

  urn: (object) => {
    if (!object) return null
    if (object.urn) return object.urn
    if (object.bsn) return `urn:oid:2.16.840.1.113883.2.4.6.3:${object.bsn}`
    if (object.agb) return `urn:oid:2.16.840.1.113883.2.4.6.1:${object.agb}`
    if (object.identifier) return object.identifier
    return null
  }
}
