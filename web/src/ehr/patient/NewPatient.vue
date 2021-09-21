<template>
  <div class="px-10 py-5">
    <h1 class="mb-4">New Patient</h1>

    <div class="p-3 bg-red-100 rounded-md" v-if="formErrors.length">
      <b>Please correct the following error(s):</b>
      <ul>
        <li v-for="error in formErrors">* {{ error }}</li>
      </ul>
    </div>

    <patient-form :value="patient" @input="(newPatient)=> {patient = newPatient}"/>

    <div class="mt-5">
      <button @click="checkForm"
              class="btn btn-primary mr-4"
      >Add Patient
      </button>

      <button type="button"
              class="btn btn-secondary"
              @click="cancel"
      >
        Cancel
      </button>
    </div>
  </div>
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
          .catch(error => this.$status.error(error))
    }
  }
}
</script>
