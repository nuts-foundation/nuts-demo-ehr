export default {
  render: (patient, organisation) => {
    update(document.getElementById('patient-pane'), patient.id, organisation);
    return Promise.resolve()
  }
}

function update (element, patient, organisation) {
  element.classList.add('loading')
  fetch(`/api/observation/remoteByPatientId/${patient}/${organisation.identifier}`)
    .then(response => {
      element.classList.remove('loading')
      if (response.status == 401) {
        // We're not authenticated (anymore), go to IRMA flow
        window.localStorage.setItem('afterLoginReturnUrl', `private/patient/${patient}/external/${organisation.identifier}`)
        window.location.hash = 'private/escalate'
        return response.text().then(t => Promise.reject(t))
      } else {
        if (response.ok) {
          return response.json()
        }
        // extract error message from body and wrap it in a rejected promise
        return response.text().then((text) => new Promise((resolve, reject) => reject(text)))
      }
    })
    .then(json => {
      element.innerHTML = observationsTemplate(json)
      addReloadListener(element, patient, organisation)
    })
    .catch(error => {
      element.innerHTML = `
        &nbsp;
        <p style="text-align: right"><button class="btn btn-primary remote-observations-reload">🔄 Reload</button></p>
        <h2>Error</h2>
        <p>Could not load remote observations: ${error}</p>
      `
      addReloadListener(element, patient, organisation)
    })
}

function addReloadListener (element, patient, organisation) {
  element.querySelector('.remote-observations-reload').addEventListener('click', () => {
    update(element, patient, organisation)
  })
}

const observationsTemplate = (observations) => `
  &nbsp;
  <p style="text-align: right"><button class="btn btn-primary remote-observations-reload">🔄 Reload</button></p>

  ${observations.length === 0 ? '<p><em>No remote observations found</em></p>' : ''}

  ${observations.map(observation => `
    <div class="card"><div class="card-body">
      <code>${observation.timestamp}</code>
      <p>${observation.content.replace('\n', '</p><p>')}</p>
    </div></div>
    &nbsp;
  `).join('')}
`
