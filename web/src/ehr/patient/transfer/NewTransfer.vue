<template>
  <div>
    <transfer-form :transfer="transfer"></transfer-form>

    <div class="mt-4">
      <button @click="createDossierAndTransfer" class="btn btn-primary">Create Transfer</button>
      <button
          class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"
          @click="cancel">Cancel</button>
    </div>
  </div>
</template>
<script>
import TransferForm from "./TransferForm.vue"

export default {
  components: {TransferForm},
  data() {
    return {
      transfer: {
        id: undefined,
        carePlan: {
          patientProblems: []
        }
      },
    }
  },
  methods: {
    cancel() {
      this.$router.push({name: 'ehr.patient', params: {id: this.$route.params.id}})
    },
    createDossierAndTransfer() {
      this.$api.createDossier({body: {patientID: this.$route.params.id, name: 'Transfer'}})
          .then(dossier => this.createTransfer(dossier.id))
          .then(transfer => this.$router.push({name: 'ehr.patient.transfer.edit', params: {transferID: transfer.id}}))
          .catch(error => this.$status.error(error))
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
