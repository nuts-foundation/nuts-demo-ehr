<template>
  <div>
    <p v-if="apiError" class="p-3 bg-red-100 rounded-md">Error: {{ apiError }}</p>

    <transfer-form :patient="patient" :transfer="transfer" @input="(newTransfer) => {this.transfer = newTransfer}"/>

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
      apiError: null,
      patient: {},
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
          .then(transfer => this.$router.push({name: 'ehr.transfer.edit', params: {id: transfer.id}}))
          .catch(error => this.apiError = error)
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
          .catch(error => this.apiError = error)
    }
  },
  mounted() {
    this.fetchPatient(this.$route.params.id)
  }
}
</script>
