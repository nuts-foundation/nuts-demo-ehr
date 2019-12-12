export default {
  // Fetch all patients and render to the table
  render: () => {
    return fetch('/api/patient/all')
    .then(response => response.json())
    .then(json => {
      document.getElementById('patient-overview').innerHTML = template(json);
    });
  }
}

const template = (patients) => `
  <h2>Patients in care</h2>

  <table class="table table-borderless table-bordered table-hover">

    <thead class="thead-dark">
      <tr>
        <th>Last name</th>
        <th>First name</th>
      </tr>
    </thead>

    <tbody>
      ${patients.map(patient => `
        <tr>
          <td><a href='#patient-details/${patient.id}'>${patient.name.family}</a></td>
          <td>${patient.name.given}</td>
        </tr>
      `).join('')}
    </tbody>

  </table>
`;
