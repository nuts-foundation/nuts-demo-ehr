const config        = require('../../util/config');
const openAPIHelper = require('./open-api-helper');

const call = openAPIHelper.call({
  baseURL:    `http://${config.nuts.node}`,
  definition: 'https://github.com/nuts-foundation/nuts-event-octopus/raw/master/docs/_static/nuts-event-store.yaml'
});

module.exports = {
  allEvents: async ()      => call('list'),
  getEvent:  async (jobId) => call('getEvent', jobId)
};
