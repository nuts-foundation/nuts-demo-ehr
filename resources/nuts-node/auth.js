const config = require('../../util/config')
var FormData = require('form-data')
const apiHelper = require('./open-api-helper')
const specExperimental = 'https://raw.githubusercontent.com/nuts-foundation/nuts-node/master/docs/_static/auth/experimental.yaml'
const specV0 = 'https://raw.githubusercontent.com/nuts-foundation/nuts-node/master/docs/_static/auth/v0.yaml'

const callExperimental = apiHelper.call({
  baseURL: `http://${config.nuts.node}`,
  definition: specExperimental
})
const callV0 = apiHelper.call({
  baseURL: `http://${config.nuts.node}`,
  definition: specV0
})

const auth = {
  drawUpContract: async (organisationDID) => callExperimental('drawUpContract', null, {
    type: 'BehandelaarLogin',
    language: 'NL',
    version: 'v3',
    legalEntity: organisationDID,
  }),
  createLoginSession: async (contract) => callExperimental('createSignSession', null, {means: 'irma', payload: contract.message, params: {}}),
  createSession: async (contract) => callV0('createSession', null, contract),
  sessionRequestStatus: async (id) => callExperimental('getSignSessionStatus', id),
  verifySignature: async (vp) => callExperimental('verifySignature', null, {VerifiablePresentation: vp}),

  validateContract: async (contract) => callV0('validateContract', null, contract),
  createJwtBearerToken: async (context) => callV0('createJwtBearerToken', null, context),
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
      definition: specV0
    })
    return otherAuth('createAccessToken', null, formData, headers)
  },

  introspectAccessToken: async (accessToken) => {
    const formData = new FormData()
    formData.append('token', accessToken)
    const headers = {
      ...formData.getHeaders()
    }
    return callV0('introspectAccessToken', null, formData, headers)
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
