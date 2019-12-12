const config   = require('../../util/config');
const patients = config.patients;

module.exports = {
  all: () => {
    return Promise.resolve(patients);
  },

  byId: (id) => {
    return Promise.resolve(patients.find((patient) => patient.id == id))
  },

  byBSN: (bsn) => {
    return Promise.resolve(patients.find((patient) => patient.bsn == bsn))
  }
};
