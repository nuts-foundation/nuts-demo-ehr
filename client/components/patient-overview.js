export default {
  // Fetch all patients and render to the table
  render: () => {
    return fetch('/api/patient/all')
      .then(response => response.json())
      .then(json => {
        document.getElementById('patient-overview').innerHTML = template(json)
      })
  }
}

const template = (patients) => `
  <h2>Patients in care</h2>

  <table class="table table-borderless table-bordered table-hover">

    <thead class="thead-dark">
      <tr>
        <th>Name</th>
        <th>Date of birth</th>
      </tr>
    </thead>

    <tbody>
      ${patients.map(patient => `
        <tr>
          <td><a href='#patient-details/${patient.id}'>${patient.name.given} ${patient.name.family}</a></td>
          <td>${patient.birthDate}</td>
        </tr>
      `).join('')}
    </tbody>

  </table>
`
