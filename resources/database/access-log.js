const events = require('../../util/events');
const logs   = [];

module.exports = {
  byPatientId: (id) => {
    return Promise.resolve(logs.filter(o => o.patientId == id));
  },

  store: (log) => {
    log.id = Math.max(...logs.map(o => o.id)) + 1;
    logs.push(log);
    events.accessLog.emit('stored', log);
    return Promise.resolve(log);
  }
};
