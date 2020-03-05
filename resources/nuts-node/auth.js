const config = require('../../util/config')
var FormData = require('form-data')
const apiHelper = require('./open-api-helper')
const definitionLocation = 'https://raw.githubusercontent.com/nuts-foundation/nuts-auth/0cdd505d38ac3387062437ffe25863ca2cc4a11a/docs/_static/nuts-auth.yaml'
const call = apiHelper.call({
  baseURL: `http://${config.nuts.node}`,
  definition: definitionLocation
})

module.exports = {
  createLoginSession: async () => await call('createSession', null, loginContract()),
  createSession: async (contract) => await call('createSession', null, contract),
  sessionRequestStatus: async (id) => await call('sessionRequestStatus', id),
  validateContract: async (contract) => await call('validateContract', null, contract),
  createJwtBearerToken: async (context) => await call('createJwtBearerToken', null, context),
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
    return await otherAuth('createAccessToken', null, formData, headers)
  },

  introspectAccessToken: async (accessToken) => {
    const formData = new FormData()
    formData.append('token', accessToken)
    const headers = {
      ...formData.getHeaders()
    }
    return await call('introspectAccessToken', null, formData, headers)
  }
}

function loginContract () {
  return {
    type: 'BehandelaarLogin',
    language: 'NL',
    version: 'v1',
    legalEntity: `urn:oid:2.16.840.1.113883.2.4.6.1:${config.organisation.agb}`
  }
}
