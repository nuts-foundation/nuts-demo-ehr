<template>
  <h1 class="mt-10 mb-6">Edit Patient</h1>

  <form @submit.stop.prevent="submit">
    <div class="p-3 bg-red-100 rounded-md" v-if="formErrors.length">
      <b>Please correct the following error(s):</b>
      <ul>
        <li v-for="error in formErrors">* {{ error }}</li>
      </ul>
    </div>


    <div class="bg-white p-5 rounded-lg shadow-lg">
      <patient-form :value="patient" mode="edit" @input="(newPatient)=> {patient = newPatient}"/>
    </div>

    <div class="mt-4">
      <button type="submit"
              :class="{'btn-loading': loading}"
              class="btn btn-primary mr-4"
      >Update Patient
      </button>

      <button type="button"
              class="btn btn-secondary"
              @click="cancel"
      >
        Cancel
      </button>
    </div>
  </form>
</template>
<script>

import PatientForm from "./PatientFields.vue";

export default {
  components: {PatientForm},
  data() {
    return {
      loading: false,
      formErrors: [],
      patient: {
        ObjectID: '',
        id: '',
        ssn: '',
        dob: null,
        gender: 'unknown',
        firstName: '',
        surname: '',
        zipcode: '',
        email: null,
      }
    }
  },
  emits: ["statusUpdate"],
  mounted() {
    this.fetchPatient()
  },
  methods: {
    cancel() {
      this.$router.push({name: 'ehr.patient', params: {id: this.patient.ObjectID}})
    },
    checkForm(e) {
      // reset the errors
      // no checks yet
      this.formErrors.length = 0
      return true
    },
    submit() {
      if (!this.checkForm()) {
        return false
      }
      this.loading = true;

      let patientID = this.$route.params.id

      this.$api.updatePatient({patientID: patientID, body: this.patient})
          .then(() => {
            this.$emit("statusUpdate", "Patient updated")
            this.$router.push({name: 'ehr.patient', params: {id: this.patient.ObjectID}})
          })
          .catch(error => this.$status.error(error))
          .finally(() => this.loading = false)
    },
    fetchPatient() {
      this.$api.getPatient({patientID: this.$route.params.id})
          .then(patient => this.patient = patient)
          .catch(error => this.$status.error(error))
    }
  },

}
</script>
