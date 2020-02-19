let patientId;

export default {
  render: (patient) => {
    patientId = patient.id;
    renderLogs(patientId);
  }
}

function renderLogs(patientId) {
  return fetch(`/api/accessLog/byPatientId/${patientId}`)
  .then(response => response.json())
  .then(logs => {
    document.getElementById('patient-logs').innerHTML = template(logs);
  });
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
          <td>AGB code ${log.user['irma-demo.nuts.agb.agbcode']}</td>
        </tr>
      `).join('') : '<tr><td colspan="3" style="text-align: center;"><em>None</em></td></tr>'}
    </tbody>
  </table>
`;
