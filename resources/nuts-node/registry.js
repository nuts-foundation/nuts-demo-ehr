const config = require('../../util/config');
const call   = require('./open-api-helper').call({
  baseURL:    `http://${config.nuts.node}`,
  definition: 'https://github.com/nuts-foundation/nuts-registry/raw/master/docs/_static/nuts-registry.yaml'
});

module.exports = {
  searchOrganizations:    async (query) => (await call('searchOrganizations', {query})).sort((a,b) => a.name.localeCompare(b.name)),
  organizationById:       async (id)    => await call('organizationById', id),
  deregisterOrganization: async (id)    => await call('deregisterOrganization', id),
  registerOrganization:   async (org)   => await call('registerOrganization', null, org)
};
