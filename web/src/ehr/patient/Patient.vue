<template>
  <div>
    <p v-if="!!error" class="m-4">Error: {{ error }}</p>

    <patient-details :patient="patient"/>

    <div class="mt-8">
      <div>
        <h1 class="text-xl float-left">Dossiers</h1>
        <button class="float-right inline-flex items-center" @click="$router.push({name: 'ehr.patient.dossier.new'})">
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
import PatientDetails from "./PatientDetails.vue";
import ModalWindow from "../../components/ModalWindow.vue";

export default {
  components: {PatientDetails, ModalWindow},
  data() {
    return {
      patient: {},
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
      this.api.getPatient({patientID: this.$route.params.id})
          .then(patient => this.patient = patient)
          .catch(reason => console.log(reason))
    }

  },
  mounted() {
    this.fetchPatient()
  },
}
</script>