<template>
  <h1>New Patient</h1>
  <patient-form :value="patient" @input="(newPatient)=> {patient = newPatient}"/>
  <button @click="checkForm"
          class="bg-blue-400 hover:bg-blue-500 text-white font-medium rounded-md px-3 py-2"
  >Update Patient
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
      apiError: '',
      formErrors: [],
      patient: {
        PatientID: '',
        id: '',
        ssn: '',
        dob: null,
        gender: 'unknown',
        firstName: 'Henk',
        surname: 'de Vries',
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
      this.$router.push({name: 'ehr.patient', params: {id: this.patient.PatientID}})
    },
    checkForm(e) {
      // reset the errors
      this.formErrors.length = 0
      this.apiError = ''

      return this.updatePatient()

      // e.preventDefault()
    },
    updatePatient() {
      let patientID = this.$route.params.id
      this.$api.put(`web/private/patient/${patientID}`, this.patient)
          .then(response => {
            this.$emit("statusUpdate", "Patient updated")
            this.$router.push({name: 'ehr.patient', params: {id: this.patient.PatientID}})
          })
          .catch(response => this.apiError = response.statusText)
    },
    fetchPatient() {
      let patientID = this.$route.params.id
      this.$api.get(`/web/private/patient/${patientID}`)
          .then(patient => this.patient = patient)
          .catch(reason => console.log(reason))
    }
  },

}
</script>