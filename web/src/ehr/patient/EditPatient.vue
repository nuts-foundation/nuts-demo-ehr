<template>
  <h1 class="mt-10 mb-6">Edit Patient</h1>

  <form @submit.stop.prevent="submit">
    <div class="bg-white rounded-lg shadow-lg">
      <div class="sticky top-0 z-10 p-3 bg-red-100 text-red-500 rounded-t-md" v-if="formErrors.length">
        <label class="text-red-500">Please correct the following error{{ formErrors.length === 0 ? '' : 's' }}:</label>

        <ul class="text-sm">
          <li v-for="error in formErrors">— {{ error }}</li>
        </ul>
      </div>

      <div class="p-5">
        <patient-form :value="patient" mode="edit" @input="(newPatient)=> {patient = newPatient}"/>
      </div>
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
  mounted() {
    this.fetchPatient()
  },
  methods: {
    cancel() {
      this.$router.push({name: 'ehr.patient', params: {id: this.patient.ObjectID}})
    },
    checkForm(e) {
      // reset the errors
      this.formErrors.length = 0

      if (this.patient.firstName === "" || this.patient.surname === "") {
        this.formErrors.push("The firstname and surname are required fields.")
      }

      if (this.patient.dob === null) {
        this.formErrors.push("The date of birth is a required field.")
      }

      return this.formErrors.length === 0
    },
    submit() {
      if (!this.checkForm()) {
        return false
      }
      this.loading = true;

      let patientID = this.$route.params.id

      this.$api.updatePatient({patientID: patientID}, this.patient)
          .then(() => {
            this.$store.commit("statusUpdate", "Patient updated")
            this.$router.push({name: 'ehr.patient', params: {id: this.patient.ObjectID}})
          })
          .catch(error => this.$status.error(error))
          .finally(() => this.loading = false)
    },
    fetchPatient() {
      this.$api.getPatient({patientID: this.$route.params.id})
          .then(result => this.patient = result.data)
          .catch(error => this.$status.error(error))
    }
  },

}
</script>
