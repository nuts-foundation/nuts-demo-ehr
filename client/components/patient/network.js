import Thimbleful from 'thimbleful';
import call       from '../../component-loader';

let patientId;

export default {
  render: async (patient) => {
    patientId = patient.id;
    stopIntervals();
    await renderNetwork();

    document.getElementById('patient-add-consent-org').addEventListener('keyup', (e) => {
      const input = e.target.value;
      if ( input.length >= 2 ) {
        fetch(`/api/organisation/search/${input}`)
        .then(response => response.json())
        .then(result => renderAutoComplete(result));
      }
      if ( input.length == 0 ) {
        renderAutoComplete([]);
      }
    });
  }
}

let receivedInterval, givenInterval;

function stopIntervals() {
  if ( receivedInterval ) {
    window.clearInterval(receivedInterval);
    receivedInterval = null;
  }
  if ( givenInterval ) {
    window.clearInterval(givenInterval);
    givenInterval = null;
  }
}

async function renderNetwork() {
  const network = document.getElementById('patient-network');
  network.innerHTML = template([], []);
  const received = network.querySelector('#received-list');
  const given    = network.querySelector('#given-list');

  if ( !receivedInterval )
    receivedInterval = window.setInterval(() => renderReceived(received), 3000);
  renderReceived(received);

  if ( !givenInterval )
    givenInterval = window.setInterval(() => renderGiven(given), 3000);
  renderGiven(given);
}

async function renderReceived(element) {
  return call(`/api/consent/${patientId}/received`, element)
  .then(json => {
    if ( json.length == 0 )
      element.innerHTML = `<li><em>None</em></li>`;
    else
      element.innerHTML = json.map(c => `<li><a href="#patient-network/${patientId}/${c.identifier}">${c.name}</a></li>`);
  })
  .catch(error => {
    element.innerHTML = `<li><em>Could not load organisations: ${error}</em></li>`;
  });
}

async function renderGiven(element) {
  return call(`/api/consent/${patientId}/given`, element)
  .then(json => {
    if ( json.length == 0 )
      element.innerHTML = `<li><em>None</em></li>`;
    else
      element.innerHTML = json.map(c => `<li>${c.name}</li>`);
  })
  .catch(error => {
    element.innerHTML = `<li><em>Could not load organisations: ${error}</em></li>`;
  });
}

function renderAutoComplete(results) {
  document.getElementById('patient-consent-auto-complete')
          .innerHTML = results.map(result => `
    <a class="list-group-item list-group-item-action d-flex justify-content-between align-items-center"
          data-organisation-id="${result.identifier}" data-organisation-name="${result.name}">
      ${result.name}
      <span class="badge badge-primary badge-pill">${result.identifier.split(':').pop()}</span>
    </a>
  `).join('');
}

// Selecting an option from the auto-complete dropdown
Thimbleful.Click.instance().register('a[data-organisation-id]', (e) => {
  const id = e.target.attributes['data-organisation-id'].value;
  const name = e.target.attributes['data-organisation-name'].value;
  document.querySelector('input[name="organisation-id"]').value = id;
  document.getElementById('patient-add-consent-org').value = name;
  document.getElementById('patient-consent-auto-complete').innerHTML = '';
});

// Storing new consent
Thimbleful.Click.instance().register('#patient-add-consent-button', (e) => {
  const organisationURN = document.querySelector('input[name="organisation-id"]').value;
  const reason = document.getElementById('patient-add-consent-reason').value;

  storeConsent({ organisationURN, reason })
  .then(() => renderNetwork());
});

function storeConsent(consent) {
  return fetch(`/api/consent/${patientId}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(consent)
  })
  .then(response => {
    if ( response.status != 201 ) throw 'Error storing observation';
  });
}

const template = (receivedConsents, givenConsents) => `
  &nbsp;

  <div class="card">
    <div class="card-body">
      <p>Organisations that have shared information with you:</p>
      <ul id='received-list'>
      ${receivedConsents.length > 0 ? receivedConsents.map(consent => `
        <li><a href="#patient-network/${patientId}/${consent.identifier}">${consent.name}</a></li>
      `).join('') : '<li><i>None</i></li>'}
      </ul>
    </div>
  </div>

  &nbsp;

  <div class="card">
    <div class="card-body">
      <p>Organisations you're sharing information with:</p>
      <ul id='given-list'>
      ${givenConsents.length > 0 ? givenConsents.map(consent => `
        <li>${consent.name}</li>
      `).join('') : '<li><i>None</i></li>'}
      </ul>

      <p><button class="btn btn-primary" data-toggle="#patient-add-consent">Add</button></p>

      <section id="patient-add-consent" class="page">
        <div class="card">
          <div class="card-body">
            <form id="patient-consent-form">

              <p>Share your information about this patient with another organisation:
              <div class="form-group row">
                <label for="patient-add-consent-org" class="col-sm-3 col-form-label">Organisation:</label>
                <div class="col-sm-9">
                  <input type="hidden" name="organisation-id"/>
                  <input type="text" class="form-control" id="patient-add-consent-org" placeholder="Organisation name" autocomplete="off"/>
                  <div class="list-group auto-complete" id="patient-consent-auto-complete"></div>
                </div>
              </div>
              <div class="form-group row">
                <label for="patient-add-consent-reason" class="col-sm-3 col-form-label">Legal basis:</label>
                <div class="col-sm-9">
                  <input type="text" name="reason" class="form-control" id="patient-add-consent-reason" placeholder="Your legal basis for sharing this information"/>
                </div>
              </div>
              <button id="patient-add-consent-button" type="button" class="btn btn-primary float-right">Share</button>

            </form>
          </div>
        </div>
      </section>
    </div>
  </div>
`;
