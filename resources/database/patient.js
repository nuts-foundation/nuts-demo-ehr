const config = require('../../util/config')
const patients = config.patients

module.exports = {
  all: () => {
    return Promise.resolve(patients)
  },

  byId: (id) => {
    const patient = patients.find((patient) => patient.id == id)
    if (!patient) {
      return Promise.reject('not found')
    }
    return Promise.resolve(patient)
  },

  byBSN: (bsn) => {
    const patient = patients.find((patient) => patient.bsn == bsn)
    if (!patient) {
      return Promise.reject('not found')
    }
    return Promise.resolve(patient)
  }
}
