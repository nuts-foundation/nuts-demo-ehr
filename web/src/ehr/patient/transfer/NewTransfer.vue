<template>
  <div>
    <transfer-form :transfer="transfer"></transfer-form>

    <div class="mt-4">
      <button @click="createDossierAndTransfer" class="btn btn-primary">Create Transfer</button>
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
        transferDate: "",
        description: "",
      },
    }
  },
  methods: {
    createDossierAndTransfer() {
      this.$api.createDossier({body: {patientID: this.$route.params.id, name: 'Transfer'}})
          .then(dossier => this.createTransfer(dossier.id))
          .then(transfer => this.$router.push({name: 'ehr.patient.transfer.edit', params: {transferID: transfer.id}}))
          .catch(error => this.$errors.report(error))
    },
    createTransfer(dossierID) {
      return this.$api.createTransfer({
        body: {
          dossierID: dossierID,
          transferDate: this.transfer.transferDate,
          description: this.transfer.description,
        }
      })
    },
    fetchPatient(patientID) {
      this.$api.getPatient({patientID: patientID})
          .then(patient => this.patient = patient)
          .catch(error => this.$errors.report(error))
    }
  },
  mounted() {
    this.fetchPatient(this.$route.params.id)
  }
}
</script>
