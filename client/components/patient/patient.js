import Thimbleful from 'thimbleful';
import Details from './details'
import Observations from './observations'
import Network from './network'
import Logs from './access-logs'

const router = new Thimbleful.Router();
let currentPatient;

router.addRoutes({
  details: () => {
    Details.render(currentPatient);
    openTab('patient-details');
  },
  observations: () => {
    Observations.render(currentPatient);
    openTab('patient-observations');
  },
  network: () => {
    Network.render(currentPatient);
    openTab('patient-network');
  },
  logs: () => {
    Logs.render(currentPatient);
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
      <li class="breadcrumb-item"><a href="#dashboard">Patients in care</a></li>
      <li class="breadcrumb-item active" aria-current="page">${patient.name.given} ${patient.name.family}</li>
    </ol>
  </nav>

  <ul class="nav nav-tabs">
    <li class="nav-item">
      <a class="nav-link active" data-open="#patient-details">Details</a>
    </li>
    <li class="nav-item">
      <a class="nav-link" data-open="#patient-observations">Observations</a>
    </li>
    <li class="nav-item">
      <a class="nav-link" data-open="#patient-logs">Logs</a>
    </li>
    <li class="nav-item">
      <a class="nav-link" data-open="#patient-network">Network</a>
    </li>
  </ul>

  <section class="tab-pane active" id="patient-details" data-group="patient-tab-panes" data-follower="a[data-open='#patient-details']"></section>
  <section class="tab-pane" id="patient-observations" data-group="patient-tab-panes" data-follower="a[data-open='#patient-observations']"></section>
  <section class="tab-pane" id="patient-logs" data-group="patient-tab-panes" data-follower="a[data-open='#patient-logs']"></section>
  <section class="tab-pane" id="patient-network" data-group="patient-tab-panes" data-follower="a[data-open='#patient-network']"></section>
`
