<template>
  <form @submit.stop.prevent="createDossierAndTransfer" novalidate>
    <div class="sticky top-0 z-10 p-3 bg-red-100 text-red-500 rounded-md shadow-sm" v-if="formErrors.length">
      <label class="text-red-500">Please correct the following error{{formErrors.length === 0 ? '' : 's'}}:</label>

      <ul class="text-sm">
        <li v-for="error in formErrors">— {{ error }}</li>
      </ul>
    </div>

    <transfer-form :transfer="transfer"></transfer-form>

    <div class="mt-6">
      <button type="submit"
              class="btn btn-primary mr-4"
              :class="{'btn-loading': loading}">
        Create Transfer
      </button>

      <button class="btn btn-secondary"
              @click="cancel">
        Cancel
      </button>
    </div>
  </form>
</template>
<script>
import TransferForm from "./TransferFields.vue"

export default {
  components: {TransferForm},
  data() {
    return {
      loading: false,
      formErrors: [],
      transfer: {
        id: undefined,
        transferDate: null,
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
    checkForm() {
      // Reset the formErrors array:
      this.formErrors.length = 0

      if (this.transfer.transferDate === null) {
        this.formErrors.push("The transfer date is a required field.")
      }

      let problemsWithInterventionsAndOrName = this.transfer.carePlan.patientProblems.filter(patientProblem => {
        return patientProblem.problem.name !== "" || !!patientProblem.interventions.find(intervention => intervention.comment !== "")
      })
      if (problemsWithInterventionsAndOrName.length === 0) {
        this.formErrors.push("The care plan should have at least 1 problem.")
      }

      let problemsWithInterventions = problemsWithInterventionsAndOrName.filter((patientProblem, i) => {
        return !!patientProblem.interventions.find(intervention => intervention.comment !== "")
      })
      if (!!problemsWithInterventions.find(patientProblem => patientProblem.problem.name === "")) {
        this.formErrors.push("Each problem must have a description")
      }

      let problemsWithoutInterventions = problemsWithInterventionsAndOrName.filter((patientProblem => !patientProblem.interventions.find(intervention => intervention.comment !== "")))
      if (problemsWithoutInterventions.length >0) {
        this.formErrors.push("Each problem must have at least 1 intervention.")
      }

      return this.formErrors.length === 0
    },
    cancel() {
      this.$router.push({name: 'ehr.patient', params: {id: this.$route.params.id}})
    },
    createDossierAndTransfer() {
      if (!this.checkForm()) {
        return
      }
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
