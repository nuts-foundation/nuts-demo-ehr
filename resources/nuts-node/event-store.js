const config = require('../../util/config')
const openAPIHelper = require('./open-api-helper')

const call = openAPIHelper.call({
  baseURL: `http://${config.nuts.node}`,
  definition: 'https://github.com/nuts-foundation/nuts-event-octopus/raw/master/docs/_static/nuts-event-store.yaml'
})

module.exports = {
  allEvents: async () => decodePayloads(await call('list')),
  getEvent: async (jobId) => await call('getEvent', jobId)
}

function decodePayloads (events) {
  for (const event of events.events || []) {
    event.payload = JSON.parse(Buffer.from(event.payload, 'base64').toString())
  }
  return events
}
