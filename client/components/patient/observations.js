import Thimbleful from 'thimbleful';

let patientId;

export default {
  render: (patient) => {
    patientId = patient.id;
    renderObservations(patientId);
  }
}

function renderObservations(patientId) {
  return fetch(`/api/observation/byPatientId/${patientId}`)
  .then(response => response.json())
  .then(observations => {
    document.getElementById('patient-observations').innerHTML = template(observations);
  });
}

const template = (observations) => `
  &nbsp;

  ${observations.map(observation => `
    <div class="card"><div class="card-body">
      <code>${observation.timestamp}</code>
      <p>${observation.content}</p>
    </div></div>
    &nbsp;
  `).join('')}

  &nbsp;

  <h4>New observation</h4>

  &nbsp;

  <textarea class="form-control" rows="5" id="new-observation"></textarea>
  <button class="btn btn-primary float-right" id="add-observation">Save</button>
`;

// Add click handler for storing new observations
Thimbleful.Click.instance().register('button#add-observation', (e) => {
  const ta = document.getElementById('new-observation');

  storeObservation({
    patientId: patientId,
    content: ta.value
  })
  .then(() => {
    ta.value = '';
    renderObservations(patientId);
  });
});

function storeObservation(observation) {
  return fetch('/api/observation', {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(observation)
  })
  .then(response => {
    if ( response.status != 201 ) throw 'Error storing observation';
  });
}
