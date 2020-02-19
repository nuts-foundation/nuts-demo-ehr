import call from '../../../component-loader';

export default {
  render: (element, patient, organisation) => {
    console.log("Rendering remote observations", Date.now());
    console.log(element);
    update(element, patient, organisation);
    return Promise.resolve();
  }
}

function update(element, patient, organisation) {
  element.classList.add('loading');
  fetch(`/api/observation/remoteByPatientId/${patient}/${organisation.identifier}`)
  .then(response => {
    element.classList.remove('loading');
    if ( response.status == 401 ) {
      // We're not authenticated (anymore), go to IRMA flow
      window.location.hash = 'irma-login';
      return response.text().then(t => Promise.reject(t));
    } else {
      return response.json();
    }
  })
  .then(json => {
    element.innerHTML = observationsTemplate(json);
    addReloadListener(element, patient, organisation);
  })
  .catch(error => {
    element.innerHTML = `
      &nbsp;
      <p style="text-align: right"><button id="remote-observations-reload">🔁 Reload</button></p>
      <h2>Error</h2><p>Could not load remote observations: ${error}</p>
    `;
    addReloadListener(element, patient, organisation);
  });
}

function addReloadListener(element, patient, organisation) {
  document.getElementById('remote-observations-reload').addEventListener('click', () => {
    update(element, patient, organisation);
  });
}

const observationsTemplate = (observations) => `
  &nbsp;
  <p style="text-align: right"><button id="remote-observations-reload">🔁 Reload</button></p>

  ${observations.length === 0 ? "<p><em>No remote observations found</em></p>" : ""}

  ${observations.map(observation => `
    <div class="card"><div class="card-body">
      <code>${observation.timestamp}</code>
      <p>${observation.content.replace('\n', '</p><p>')}</p>
    </div></div>
    &nbsp;
  `).join('')}
`;
