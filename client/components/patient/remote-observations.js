let patientId;

export default {
  render: (patient, organisationURN) => {
    patientId = patient.id;
    renderObservations(patientId)
  }
}

async function renderObservations(patientId) {
  return fetch(`/api/observation/remoteByPatientId/${patientId}`)
  .then(response => response.json())
  .then(observations => {
    document.getElementById('remote-patient-observations').innerHTML = template(observations);
  });
}

const template = (observations) => `
  &nbsp;

  ${!observations.length ? "No remote observations found" : ""}

  ${observations.map(observation => `
    <div class="card"><div class="card-body">
      <code>${observation.timestamp}</code>
      <p>${observation.content.replace('\n', '</p><p>')}</p>
    </div></div>
    &nbsp;
  `).join('')}

`;
