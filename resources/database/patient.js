const config = require('../../util/config')
const patients = config.patients

module.exports = {
  all: () => {
    return Promise.resolve(patients)
  },

  byId: (id) => {
    const patient = patients.find((patient) => patient.id === id)
    if (!patient) {
      throw new Error('not found')
    }
    return Promise.resolve(patient)
  },

  byBSN: (bsn) => {
    const patient = patients.find((patient) => patient.bsn === bsn)
    if (!patient) {
      throw new Error('not found')
    }
    return Promise.resolve(patient)
  }
}
