<template>
  <div class="dossier-container mt-8 flex flex-col">

    <div class="dossier-header">
      <h1 class="text-xl float-left">Dossiers</h1>
      <button class="float-right inline-flex items-center" @click="$router.push({name: 'ehr.patient.dossier.new'})">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M9 13h6m-3-3v6m5 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
        </svg>
        <span>Add Dossier</span>
      </button>
    </div>

    <div class="dossier-list">

      <table v-if="collaborationDossiers" class="min-w-full divide-y divide-gray-200">
        <thead>
        <tr>
          <th class="text-left">Name</th>
          <th class="text-left">Status</th>
          <!--          <th class="text-left">Network</th>-->
        </tr>
        </thead>
        <tbody>
        <tr class="cursor-pointer"
            @click="openDossier(dossier)"
            v-for="dossier in collaborationDossiers">
          <td>{{ dossier.name }}</td>
          <td>{{ dossier.transfer ? dossier.transfer.status : "" }}</td>
          <!--          <td>{{ dossier.network.join(', ') }}</td>-->
        </tr>
        </tbody>
      </table>

      <div v-else class="min-w-full">
        No dossiers for this patient found.
        <router-link :to="{name: 'ehr.patient.dossier.new'}">Create one!</router-link>
      </div>
    </div>

  </div>

  <div class="reports-container mt-8">

    <div class="reports-header">
      <h1 class="text-xl float-left">Reports</h1>
      <button class="float-right inline-flex items-center">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M9 13h6m-3-3v6m5 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
        </svg>
        <span>Add Report</span>
      </button>

    </div>

    <table class="reports-list min-w-full divide-y divide-gray-200">
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
</template>
<script>
export default {
  data() {
    return {
      dossiers: [],
      transfers: [],
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
    }
  },
  computed: {
    collaborationDossiers() {
      return this.dossiers.map((dossier) => {
        dossier.transfer = this.transfers.find(transfer => transfer.dossierID === dossier.id)
        return dossier
      })
    }

  },
  methods: {
    truncate(str, n) {
      return (str.length > n) ? str.substr(0, n - 1) + '...' : str
    },
    fetchDossiers() {
      this.$api.getDossier({patientID: this.$route.params.id})
          .then(dossiers => this.dossiers = dossiers)
          .catch(error => {
            this.$errors.report(error)
            console.log(error)
          })
    },
    fetchTransfers() {
      this.$api.getPatientTransfers({patientID: this.$route.params.id})
          .then(transfers => this.transfers = transfers)
          .catch(error => {
            this.$errors.report(error)
            console.log(error)
          })
    },
    openDossier(dossier) {
      const patientID = this.$route.params.id
      if (dossier.transfer) {
        this.$router.push({name: 'ehr.patient.transfer.edit', params: {id: patientID, transferID: dossier.transfer.id}})
      }
    }
  },
  mounted() {
    this.fetchDossiers()
    this.fetchTransfers()
  },
}
</script>