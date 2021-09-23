<template>
  <div>
    <transfer-form :transfer="transfer"></transfer-form>

    <div class="mt-6">
      <button @click="createDossierAndTransfer" class="btn btn-primary mr-4" :class="{'btn-loading': loading}">Create Transfer</button>

      <button
          class="btn btn-secondary"
          @click="cancel">Cancel
      </button>
    </div>
  </div>
</template>
<script>
import TransferForm from "./TransferForm.vue"

export default {
  components: {TransferForm},
  data() {
    return {
      loading: false,
      transfer: {
        id: undefined,
        carePlan: {
          patientProblems: [
            {
              problem: {name: ""},
              interventions: [{comment: ""}]
            }
          ]
        }
      },
    }
  },
  methods: {
    cancel() {
      this.$router.push({name: 'ehr.patient', params: {id: this.$route.params.id}})
    },
    createDossierAndTransfer() {
      this.loading = true;

      this.$api.createDossier({body: {patientID: this.$route.params.id, name: 'Transfer'}})
          .then(dossier => this.createTransfer(dossier.id))
          .then(transfer => this.$router.push({name: 'ehr.patient.transfer.edit', params: {transferID: transfer.id}}))
          .catch(error => this.$status.error(error))
          .finally(() => this.loading = false)
    },
    createTransfer(dossierID) {
      return this.$api.createTransfer({
        body: {
          dossierID: dossierID,
          ...this.transfer
        }
      })
    },
    fetchPatient(patientID) {
      this.$api.getPatient({patientID: patientID})
          .then(patient => this.patient = patient)
          .catch(error => this.$status.error(error))
    }
  },
  mounted() {
    this.fetchPatient(this.$route.params.id)
  }
}
</script>
