import Thimbleful from 'thimbleful';
import details from './details'
import observations from './observations'
import network from './network'
import logs from './access-logs'

const router = new Thimbleful.Router();
let currentPatient;

router.addRoutes({
  details: () => {
    details.render(currentPatient);
    openTab('patient-details');
  },
  observations: () => {
    observations.render(currentPatient);
    openTab('patient-observations');
  },
  network: () => {
    network.render(currentPatient);
    openTab('patient-network');
  },
  logs: () => {
    logs.render(currentPatient);
    openTab('patient-logs');
  }
});

function openTab(tab) {
  document.querySelector('#patient .nav li a.active').classList.remove('active')
  document.querySelector(`#patient .nav li a#${tab}`).classList.add('active')
}

export default {
  render: (patientId, subpath) => {
    return fetch(`/api/patient/${patientId}`)
      .then(response => response.json())
      .then(patient => {
        currentPatient = patient;
        document.getElementById('patient').innerHTML = template(patient);
        router.route(subpath);
      })
  },

  rendered: () => {
    return document.getElementById('patient').children.length > 0
  }
}

const template = (patient) => `
  <nav aria-label="breadcrumb">
    <ol class="breadcrumb">
      <li class="breadcrumb-item"><a href="#private/dashboard">Patients in care</a></li>
      <li class="breadcrumb-item active" aria-current="page">${patient.name.given} ${patient.name.family}</li>
    </ol>
  </nav>

  <ul class="nav nav-tabs">
    <li class="nav-item">
      <a class="nav-link active" id="patient-details" href="#private/patient/${patient.id}/details">Details</a>
    </li>
    <li class="nav-item">
      <a class="nav-link" id="patient-observations" href="#private/patient/${patient.id}/observations">Observations</a>
    </li>
    <li class="nav-item">
      <a class="nav-link" id="patient-logs" href="#private/patient/${patient.id}/logs">Logs</a>
    </li>
    <li class="nav-item">
      <a class="nav-link" id="patient-network" href="#private/patient/${patient.id}/network">Network</a>
    </li>
  </ul>

  <section class="tab-pane active" id="patient-pane"></section>
`
