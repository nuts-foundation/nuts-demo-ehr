const config = require('../../util/config')
const openAPIHelper = require('./open-api-helper')

const call = openAPIHelper.call({
  baseURL: `http://${config.nuts.node}`,
  definition: 'https://github.com/nuts-foundation/nuts-consent-store/raw/master/docs/_static/nuts-consent-store.yaml'
})

module.exports = {
  consentsFor: async (parties) => call('queryConsent', null, query(parties)),
  checkConsent: async (parties) => call('checkConsent', null, combination(parties))
}

function query (parties) {
  return {
    custodian: openAPIHelper.urn(parties.custodian),
    actor: openAPIHelper.urn(parties.actor),
    subject: openAPIHelper.urn(parties.subject),
    page: {
      offset: 0,
      limit: 0
    },
    validAt: null
  }
}

function combination (parties) {
  return {
    ...parties,
    dataClass: 'urn:oid:1.3.6.1.4.1.54851.1:MEDICAL'
  }
}
