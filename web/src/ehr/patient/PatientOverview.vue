<template>
  <div class="mt-10 bg-white px-7 py-5 rounded-lg shadow-sm flex flex-col">
    <div class="flex justify-between items-center mb-3">
      <h2>Dossiers</h2>

      <button id="new-dossier-button"
              class="inline-flex items-center bg-nuts w-10 h-10 rounded-lg justify-center shadow-md"
              @click="$router.push({name: 'ehr.patient.dossier.new'})">
        <svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 0 24 24" width="24px" fill="#fff">
          <path d="M0 0h24v24H0V0z" fill="none"/>
          <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/>
        </svg>
      </button>
    </div>

    <div class="dossier-list">

      <table v-if="episodeDossiers.length > 0" class="min-w-full divide-y divide-gray-200">
        <thead>
        <tr>
          <th>Name</th>
          <th>Status</th>
          <th>Network</th>
        </tr>
        </thead>
        <tbody>
        <tr class="cursor-pointer"
            @click="openDossier(dossier)"
            v-for="dossier in episodeDossiers">
          <td>{{ dossier.name }}</td>
          <td>{{ dossier.transfer ? dossier.transfer.status : "" }}</td>
          <td>
            {{
              dossier.transfer && dossier.transfer.negotiations ? dossier.transfer.negotiations.map(n => n.organization.name).join(', ') : ""
            }}
          </td>
        </tr>
        </tbody>
      </table>

      <div v-else-if="loadingDossiers">
        Loading...
      </div>

      <div v-else class="min-w-full">
        No dossiers for this patient found.
        <router-link :to="{name: 'ehr.patient.dossier.new'}">Create one!</router-link>
      </div>
    </div>

  </div>

  <div class="bg-white px-7 py-5 rounded-lg shadow-sm mt-8">
    <div class="flex justify-between items-center mb-3">
      <h2>Reports</h2>
    </div>

    <table class="reports-list min-w-full divide-y divide-gray-200">
      <thead>
      <tr>
        <th>Type</th>
        <th>Value</th>
        <th>Source</th>
        <th>Episode</th>
      </tr>
      </thead>
      <tbody>
      <tr v-for="report in reports">
        <td>{{ report.type }}</td>
        <td>{{ truncate(report.value, 30) }}</td>
        <td>{{ report.source }}</td>
        <td>{{ report.episodeName }}</td>
      </tr>
      </tbody>
    </table>
  </div>
</template>
<script>
export default {
  data() {
    return {
      loadingDossiers: false,
      dossiers: [],
      transfers: [],
      reports: [],
    }
  },
  computed: {
    episodeDossiers() {
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
      this.loadingDossiers = true
      this.$api.getDossier({patientID: this.$route.params.id})
          .then(dossiers => {
            this.dossiers = dossiers
            this.loadingDossiers = false
          })
          .catch(error => {
            this.loadingDossiers = false
            this.$status.error(error)
          })
    },
    fetchReports() {
      this.$api.getReports({patientID: this.$route.params.id})
          .then(reports => this.reports = reports)
          .catch(error => this.$status.error(error))
    },
    fetchTransfers() {
      this.$api.getPatientTransfers({patientID: this.$route.params.id})
          .then(transfers => {
            this.transfers = transfers
            // Also fetch negotiations so we can show the "network" of the dossier
            return Promise.all(this.transfers.map(t => this.$api.listTransferNegotiations({transferID: t.id})))
          })
          .then(negotiations => {
            for (let i = 0; i < negotiations.length; i++) {
              this.transfers[i].negotiations = negotiations[i]
            }
          })
          .catch(error => this.$status.error(error))
    },
    openDossier(dossier) {
      const patientID = this.$route.params.id
      if (dossier.transfer) {
        this.$router.push({name: 'ehr.patient.transfer.edit', params: {id: patientID, transferID: dossier.transfer.id}})
      } else {
        this.$router.push({name: 'ehr.patient.episode.edit', params: {id: patientID, episodeID: dossier.id}})
      }
    }
  },
  mounted() {
    this.fetchDossiers()
    this.fetchReports()
    this.fetchTransfers()
  },
}
</script>
