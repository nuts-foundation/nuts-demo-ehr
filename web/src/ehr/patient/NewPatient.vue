<template>
  <div class="px-12 py-8">
    <button type="button" @click="() => this.$router.push({name: 'ehr.patients'})" class="btn btn-link mb-12">
      <span class="w-6 mr-1">
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="#000000"><path d="M0 0h24v24H0V0z"
                                                                                         fill="none"/><path
            d="M15.41 7.41L14 6l-6 6 6 6 1.41-1.41L10.83 12l4.58-4.59z"/></svg>
      </span>
      Back to patient list
    </button>

    <h1 class="mb-4">New Patient</h1>

    <form @submit.stop.prevent="confirm">
      <div class="bg-white rounded-lg shadow-lg">
        <div class="sticky top-0 z-10 p-3 bg-red-100 text-red-500 rounded-t-md" v-if="formErrors.length">
          <label class="text-red-500">Please correct the following error{{formErrors.length === 0 ? '' : 's'}}:</label>

          <ul class="text-sm">
            <li v-for="error in formErrors">— {{ error }}</li>
          </ul>
        </div>

        <div class="p-5">
          <patient-form class="space-y-5" :value="patient" @input="(newPatient)=> {patient = newPatient}"/>
        </div>
      </div>

      <div class="mt-5">
        <button type="submit"
                class="btn btn-primary mr-4"
                :class="{'btn-loading': loading}"
        >Add Patient
        </button>

        <button type="button"
                class="btn btn-secondary"
                @click="cancel"
        >
          Cancel
        </button>
      </div>
    </form>
  </div>
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
  methods: {
    cancel() {
      this.$router.push({name: 'ehr.patients'})
    },
    checkForm(e) {
      // reset the errors
      this.formErrors.length = 0

      if (this.patient.ssn === "") {
        this.formErrors.push("SSN is a required field.")
      }

      if (this.patient.firstName === "" || this.patient.surname === "") {
        this.formErrors.push("The firstname and surname are required fields.")
      }

      if (this.patient.dob === null) {
        this.formErrors.push("The date of birth is a required field.")
      }

      return this.formErrors.length === 0
    },
    confirm() {
      if (!this.checkForm()) {
        return false
      }
      this.loading = true;

      this.$api.newPatient({body: this.patient})
          .then(response => {
            this.$store.commit("statusUpdate", "Patient added")
            this.$router.push({name: 'ehr.patients'})
          })
          .catch(error => this.$status.error(error))
          .finally(() => this.loading = false)
    }
  }
}
</script>
