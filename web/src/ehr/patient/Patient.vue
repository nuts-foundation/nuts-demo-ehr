<template>
  <div>
    <patient-details :patient="patient"/>

    <router-view></router-view>
  </div>
</template>
<script>
import PatientDetails from "./PatientDetails.vue";
import ModalWindow from "../../components/ModalWindow.vue";

export default {
  components: {PatientDetails, ModalWindow},
  data() {
    return {
      patient: {},
    }
  },
  methods: {
    truncate(str, n) {
      return (str.length > n) ? str.substr(0, n - 1) + '...' : str
    },
    fetchPatient() {
      this.$api.getPatient({patientID: this.$route.params.id})
          .then(patient => this.patient = patient)
          .catch(error => this.$errors.report(error))
    }

  },
  mounted() {
    this.fetchPatient()
  },
}
</script>