const config = require('../../util/config')
const call = require('./open-api-helper').call({
  baseURL: `http://${config.nuts.node}`,
  definition: 'https://github.com/nuts-foundation/nuts-registry/raw/master/docs/_static/nuts-registry.yaml'
})

module.exports = {
  searchOrganizations: async (query) => (await call('searchOrganizations', { query })).sort((a, b) => a.name.localeCompare(b.name)),
  organizationById: async (id) => call('organizationById', { id }),
  endpointsByOrganisationId: async (id, type) => call('endpointsByOrganisationId', { orgIds: id, type }),
  deregisterOrganization: async (id) => call('deregisterOrganization', id),
  registerOrganization: async (org) => call('registerOrganization', null, org)
}
