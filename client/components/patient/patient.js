import Details      from './details';
import Observations from './observations';
import RemoteObservations from './remote-observations';
import Network      from './network';

export default {
  render: (patientId) => {
    return fetch(`/api/patient/${patientId}`)
    .then(response => response.json())
    .then(patient => {
      document.getElementById('patient').innerHTML = template(patient);

      // Render child components
      Details.render(patient);
      Observations.render(patient);
      RemoteObservations.render(patient);
      Network.render(patient);
    });
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
      <a class="nav-link" data-open="#remote-patient-observations">External Observations</a>
    </li>
    <li class="nav-item">
      <a class="nav-link" data-open="#patient-network">Network</a>
    </li>
  </ul>

  <section class="tab-pane active" id="patient-details" data-group="patient-tab-panes" data-follower="a[data-open='#patient-details']"></section>
  <section class="tab-pane" id="patient-observations" data-group="patient-tab-panes" data-follower="a[data-open='#patient-observations']"></section>
  <section class="tab-pane" id="remote-patient-observations" data-group="patient-tab-panes" data-follower="a[data-open='#remote-patient-observations']"></section>
  <section class="tab-pane" id="patient-network" data-group="patient-tab-panes" data-follower="a[data-open='#patient-network']"></section>
`;
