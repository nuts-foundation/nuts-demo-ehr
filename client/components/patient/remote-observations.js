let patientId;

export default {
  render: (patient) => {
    patientId = patient.id;
    renderButton()
  }
}

async function renderButton() {
    document.getElementById('remote-patient-observations').innerHTML = buttonTemplate();
}

async function renderObservations(patientId) {
  return fetch(`/api/observation/remoteByPatientId/${patientId}`)
  .then(response => response.json())
  .then(observations => {
    if (window.irmaLogin) {
      document.getElementById('remote-patient-observations').innerHTML = observationsTemplate(observations);
    } else {
      window.location.hash = 'irma-login'
    }
  });
}

const buttonTemplate = () => `
  &nbsp;

  <button id="load-observations">Load external observations</button>
`;


const observationsTemplate = (observations) => `

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

// Add click handler for storing new observations
Thimbleful.Click.instance().register('button#load-observations', (e) => {
  renderObservations(patientId)
});
