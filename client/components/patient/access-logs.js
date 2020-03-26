import io from '../../socketio'
const socket = io.accessLogs()

socket.on('logs', m => {
  if (!m) return
  document.getElementById('patient-logs').innerHTML = template(m)
})

export default {
  render: patient => {
    socket.emit('subscribe', patient.id)
    return Promise.resolve()
  }
}

const template = (logs) => `
  <table class="table table-borderless table-bordered table-hover">
    <thead class="thead-dark">
      <tr>
        <th>Timestamp</th>
        <th>Organisation</th>
        <th>Person</th>
      </tr>
    </thead>
    <tbody>
      ${logs.length > 0 ? logs.map(log => `
        <tr>
          <td>${new Date(log.timestamp).toLocaleString('nl-NL')}</td>
          <td>${log.actor.name}</td>
          <td>${log.user}</td>
        </tr>
      `).join('') : '<tr><td colspan="3" style="text-align: center;"><em>None</em></td></tr>'}
    </tbody>
  </table>
`
