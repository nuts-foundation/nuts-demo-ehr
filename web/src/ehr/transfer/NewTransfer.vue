<template>
  <div>
    <div class="sticky top-0 z-10 p-3 bg-red-100 rounded-md" v-if="formErrors.length">
      <b>Please correct the following error(s):</b>
      <ul>
        <li v-for="error in formErrors">* {{ error }}</li>
      </ul>
    </div>
    <transfer-form :transfer="transfer"></transfer-form>

    <div class="mt-6">
      <button @click="createDossierAndTransfer" class="btn btn-primary mr-4" :class="{'btn-loading': loading}">Create
        Transfer
      </button>

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

      if (this.transfer.carePlan.patientProblems.length === 0) {
        this.formErrors.push("The care plan should have at least 1 problem.")
      }

      this.transfer.carePlan.patientProblems.forEach((patientProblem, i) => {
        var hasInterventions = false
        patientProblem.interventions.forEach((intervention, j) => {
          // dont perform the check for the empty intervention placeholder
          if ((i > 0 || j != 0) && j == patientProblem.interventions.length - 1) {
            return
          }
          if (intervention.comment == "") {
            this.formErrors.push("Each intervention must have a comment.")
          } else {
            hasInterventions = true
          }
        })
        // Skip the check when:
        if ( // Its the placeholder, so last problem, but not the only one
            i != 0 && i == this.transfer.carePlan.patientProblems.length - 1 &&
            // and it does not have any interventions
            !hasInterventions
        ) {
          return
        }
        if (patientProblem.problem.name === "") {
          this.formErrors.push("Each problem must have a description.")
          return false
        }
        if (patientProblem.interventions.length == 0) {
          this.formErrors.push("Each problem must have at least 1 intervention.")
        }
      })
      return this.formErrors.length == 0
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
