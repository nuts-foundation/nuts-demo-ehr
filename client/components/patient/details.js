export default {
  render: (patient) => {
    document.getElementById('patient-details').innerHTML = template(patient);
  }
}

const template = (patient) => `
  <table class="table table-borderless table-bordered">
    <tbody>
      <tr>
        <th>BSN</th>
        <td>${patient.bsn}</td>
      </tr>
      <tr>
        <th>Last name</th>
        <td>${patient.name.family}</td>
      </tr>
      <tr>
        <th>First name</th>
        <td>${patient.name.given}</td>
      </tr>
      <tr>
        <th>Date of birth</th>
        <td>${patient.birthDate}</td>
      </tr>
      <tr>
        <th>Gender</th>
        <td>${patient.gender}</td>
      </tr>
    </tbody>
  </table>
`;
