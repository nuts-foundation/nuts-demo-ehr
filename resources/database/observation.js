const config = require('../../util/config')
const observations = config.observations

module.exports = {
  byPatientId: (id) => {
    return Promise.resolve(observations.filter((o) => o.patientId == id))
  },

  store: (observation) => {
    observation.id = Math.max(...observations.map(o => o.id)) + 1
    observations.push(observation)
    return Promise.resolve(observation)
  }
}
