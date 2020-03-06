const events = require('../util/events')

const {
  accessLog
} = require('../resources/database')

module.exports = async io => {
  let watchedPatients = []

  io.on('connection', socket => {
    // Subscribe to events concerning this patient
    let myWatchedPatient = null
    socket.on('subscribe', async patientId => {
      // Unsubscribe from previous patient
      watchedPatients = watchedPatients.filter(p =>
        p.patientId != myWatchedPatient ||
        p.socket !== socket
      )

      // Subscribe to new patient
      myWatchedPatient = patientId
      watchedPatients.push({ socket, patientId })

      // Send current status of this patient
      socket.emit('logs', await accessLog.byPatientId(patientId))
    })
  })

  events.accessLog.on('stored', log => {
    const notify = watchedPatients.filter(p => p.patientId == log.patientId)
    notify.forEach(async p =>
      p.socket.emit('logs', await accessLog.byPatientId(p.patientId))
    )
  })
}
