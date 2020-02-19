const logs = [];

module.exports = {
  byPatientId: (id) => {
    return Promise.resolve(logs.filter(o => o.patientId == id));
  },

  store: (log) => {
    log.id = Math.max(...logs.map(o => o.id)) + 1;
    logs.push(log);
    return Promise.resolve(log);
  }
};
