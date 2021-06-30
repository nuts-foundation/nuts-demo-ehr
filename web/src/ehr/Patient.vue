<template>
  <div>
    <p v-if="!!error" class="m-4">Error: {{ error }}</p>

    <div>
      <div class="grid grid-cols-2 gap-2 m4">
        <img v-bind:src="patient.photo" class="w-24">
        <div>
          <h1 class="text-2xl">{{ patient.surname }}, {{ patient.firstName }}</h1>
          <div class="grid grid-cols-2">
            <div>SSN: {{ patient.ssn }}</div>
            <div>Birth date: {{ patient.dob }}</div>
            <div>ID: {{ patient.id }}</div>
            <div>E-mail: {{ patient.email }}</div>
          </div>
        </div>
      </div>
    </div>

    <div class="mt-8">
      <div>
        <h1 class="text-xl float-left">Dossiers</h1>
        <button class="float-right inline-flex items-center">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M9 13h6m-3-3v6m5 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
          </svg>
          <span>Add Dossier</span>
        </button>
      </div>

      <table class="min-w-full divide-y divide-gray-200">
        <thead>
        <tr>
          <th class="text-left">Name</th>
          <th class="text-left">Status</th>
          <th class="text-left">Network</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="dossier in dossiers">
          <td>{{ dossier.name }}</td>
          <td>{{ dossier.status }}</td>
          <td>{{ dossier.network.join(', ') }}</td>
        </tr>
        </tbody>
      </table>
    </div>

    <div class="mt-8">
      <div>
        <h1 class="text-xl float-left">Reports</h1>
        <button class="float-right inline-flex items-center">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M9 13h6m-3-3v6m5 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
          </svg>
          <span>Add Report</span>
        </button>

      </div>

      <table class="min-w-full divide-y divide-gray-200">
        <thead>
        <tr>
          <th>Date</th>
          <th>Type</th>
          <th>Value</th>
          <th>Source</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="report in reports">
          <td>{{ report.date }}</td>
          <td>{{ report.type }}</td>
          <td>{{ truncate(report.value, 30) }}</td>
          <td>{{ report.source }}</td>
        </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
<script>
import patientPhoto from '../img/patients/vries.jpg';

export default {
  data() {
    return {
      patient: {
        id: "PM00567",
        ssn: 99999880,
        dob: "1981-03-01",
        firstName: "Henk",
        surname: "de Vries",
        email: "hdevries@securemail.nuts",
        photo: patientPhoto,
      },
      dossiers: [
        {
          name: "Thuiszorg",
          status: "active",
          network: ["TZ de Kastanjeboom", "HA de Leeuw"],
        },
        {
          name: "Overdracht",
          status: "requested",
          network: ["Ziekenhuis", "Thuiszorg"],
        },
        {
          name: "Overdracht",
          status: "completed",
          network: ["Ziekenhuis", "Thuiszorg"],
        },
      ],
      reports: [
        {
          date: "2021-06-01",
          type: "heartbeat",
          value: "72 bpm",
          source: "HA de Leeuw"
        },
        {
          date: "2021-06-02",
          type: "text",
          value: "Meneer is gevallen op maandag en heeft veel pijn aan de linkerheup.",
          source: "TZ de Kastanjeboom"
        }
      ],
      error: null,
    }
  },
  methods: {
    truncate(str, n) {
      return (str.length > n) ? str.substr(0, n - 1) + '...' : str
    },
    fetchPatient() {
      let patientID = this.$route.params.id
      this.$api.get(`/web/private/patient/${patientID}`)
          .then(patient => this.patient = { ...this.patient, ... patient})
          .catch(reason => console.log(reason))
    }

  },
  mounted() {
    this.fetchPatient()
  },
}
</script>