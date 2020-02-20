import Thimbleful from 'thimbleful';

import io from '../../socketio';
const socket = io.consent();

socket.on('receivedConsents', m => {
  document.querySelector('#patient-network #received-list').innerHTML =
    m.sort((a,b) => a.name.localeCompare(b.name))
     .map(c => `<li><a href="#patient-network/${patientId}/${c.identifier}">${c.name}</a></li>`)
     .join('') || '<li><i>None</i></li>';
});

socket.on('givenConsents', m => {
  document.querySelector('#patient-network #given-list').innerHTML =
    m.sort((a,b) => a.name.localeCompare(b.name))
     .map(c => `<li>${c.name}</li>`)
     .join('') || '<li><i>None</i></li>';
});

let patientId;

export default {
  render: async (patient) => {
    patientId = patient.id;
    const network = document.getElementById('patient-network');
    network.innerHTML = template();

    socket.emit('subscribe', patient.id);

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
  storeConsent({ organisationURN, reason });
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

const template = () => `
  &nbsp;

  <div class="card">
    <div class="card-body">
      <p>Organisations that have shared information with you:</p>
      <ul id='received-list'>
        <li><i>None</i></li>
      </ul>
    </div>
  </div>

  &nbsp;

  <div class="card">
    <div class="card-body">
      <p>Organisations you're sharing information with:</p>
      <ul id='given-list'>
        <li><i>None</i></li>
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
