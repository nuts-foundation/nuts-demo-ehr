<template>
  <h1 class="text-2xl">New Patient</h1>
  <div class="p-3 bg-red-100 rounded-md" v-if="formErrors.length">
    <b>Please correct the following error(s):</b>
    <ul>
      <li v-for="error in formErrors">* {{ error }}</li>
    </ul>
  </div>
  <patient-form :value="patient" @input="(newPatient)=> {patient = newPatient}"/>
  <button @click="checkForm"
      class="bg-blue-400 hover:bg-blue-500 text-white font-medium rounded-md px-3 py-2"
  >Add Patient
  </button>
  <button type="button"
          class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"
          @click="cancel"
  >
    Cancel
  </button>
</template>
<script>

import PatientForm from "./PatientForm.vue";

export default {
  components: {PatientForm},
  data() {
    return {
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
  emits: ["statusUpdate"],
  methods: {
    cancel() {
      this.$router.push({name: 'ehr.patients'})
    },
    checkForm(e) {
      // reset the errors
      this.formErrors.length = 0

      if (this.patient.ssn === "") {
        this.formErrors.push("SSN is a required field.")
        return
      }

      return this.confirm()

      // e.preventDefault()
    },
    confirm() {
      this.$api.newPatient({body: this.patient})
          .then(response => {
            this.$emit("statusUpdate", "Patient added")
            this.$router.push({name: 'ehr.patients'})
          })
          .catch(error => this.$errors.report(error))
    }
  }
}
</script>