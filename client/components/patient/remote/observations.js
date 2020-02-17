import call from '../../../component-loader';
let intervals = {};
let lastPatient;

export default {
  render: (element, patient, organisation) => {

    if ( !window.irmaLogin ) {
      return window.location.hash = 'irma-login';
    }

    if ( patient != lastPatient ) {
      Object.values(intervals).map(i => window.clearInterval(i));
      intervals = {};
      lastPatient = patient;
    }

    if ( !intervals[organisation] )
      intervals[organisation] = window.setInterval(() => update(element, patient, organisation), 3000);
    update(element, patient, organisation);

    return Promise.resolve();
  }
}

function update(element, patient, organisation) {
  call(`/api/observation/remoteByPatientId/${patient}/${organisation.identifier}`, element)
  .then(json => {
    element.innerHTML = observationsTemplate(json);
  })
  .catch(error => {
    element.innerHTML = `<h2>Error</h2><p>Could not load remote observations: ${error}</p>`;
  });
}

const observationsTemplate = (observations) => `
  &nbsp;

  ${observations.length === 0 ? "<p><em>No remote observations found</em></p>" : ""}

  ${observations.map(observation => `
    <div class="card"><div class="card-body">
      <code>${observation.timestamp}</code>
      <p>${observation.content.replace('\n', '</p><p>')}</p>
    </div></div>
    &nbsp;
  `).join('')}
`;
