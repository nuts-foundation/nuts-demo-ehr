const config = require('../../util/config')
var FormData = require('form-data')
const apiHelper = require('./open-api-helper')
const definitionLocation = 'https://raw.githubusercontent.com/nuts-foundation/nuts-auth/0cdd505d38ac3387062437ffe25863ca2cc4a11a/docs/_static/nuts-auth.yaml'
const call = apiHelper.call({
  baseURL: `http://${config.nuts.node}`,
  definition: definitionLocation
})

const auth = {
  createLoginSession: async () => call('createSession', null, loginContract()),
  createSession: async (contract) => call('createSession', null, contract),
  sessionRequestStatus: async (id) => call('sessionRequestStatus', id),
  validateContract: async (contract) => call('validateContract', null, contract),
  createJwtBearerToken: async (context) => call('createJwtBearerToken', null, context),
  createAccessToken: async (baseUrl, jwtBearerToken) => {
    // createAccessToken is a bit weird since it follows the OAuth spec and needs a FormData object instead of a plain json document
    // The baseUrl is added since the createAccessToken is usually performed on an other Nuts node than your own
    const formData = new FormData()
    formData.append('grant_type', 'urn:ietf:params:oauth:grant-type:jwt-bearer')
    formData.append('assertion', jwtBearerToken)
    const headers = {
      'X-Nuts-LegalEntity': 'Demo EHR',
      ...formData.getHeaders()
    }
    const otherAuth = apiHelper.call({
      baseURL: baseUrl,
      definition: definitionLocation
    })
    return otherAuth('createAccessToken', null, formData, headers)
  },

  introspectAccessToken: async (accessToken) => {
    const formData = new FormData()
    formData.append('token', accessToken)
    const headers = {
      ...formData.getHeaders()
    }
    return call('introspectAccessToken', null, formData, headers)
  },

  obtainAccessToken: async (context, endpoint) => {
    // Get the JWT Bearer token at the local Nuts node
    let jwtBearerTokenResponse
    try {
      jwtBearerTokenResponse = await auth.createJwtBearerToken(context)
      console.log(jwtBearerTokenResponse)
    } catch (e) {
      console.log(e)
      throw Error(`error while creating jwt bearer token: ${e.response.data}`)
    }

    const accessTokenEndpoint = endpoint.properties.authorizationServerURL

    if (!accessTokenEndpoint) {
      throw Error('no authorizationServerURL found in endpoint.properties')
    }

    console.log('accessTokenEndpoint:', accessTokenEndpoint)

    // Get the access token at the custodians Nuts node
    let accessTokenResponse
    try {
      accessTokenResponse = await auth.createAccessToken(accessTokenEndpoint, jwtBearerTokenResponse.bearer_token)
      console.log(accessTokenResponse)
    } catch (e) {
      let error
      if (e.response) {
        error = JSON.stringify(e.response.data)
      } else {
        error = e
      }
      throw Error(`error while creating access token: ${error}`)
    }

    return accessTokenResponse.access_token
  }
}

module.exports = auth

function loginContract () {
  return {
    type: 'BehandelaarLogin',
    language: 'NL',
    version: 'v1',
    legalEntity: config.nuts.version == "0.12" ? config.organisation.name : `urn:oid:2.16.840.1.113883.2.4.6.1:${config.organisation.agb}`
  }
}
