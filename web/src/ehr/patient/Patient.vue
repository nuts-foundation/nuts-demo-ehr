<template>
  <div class="px-12 py-8">
    <button type="button" @click="back" class="btn btn-link mb-12">
      <span class="w-6 mr-1">
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="#000000"><path d="M0 0h24v24H0V0z"
                                                                                         fill="none"/><path
            d="M15.41 7.41L14 6l-6 6 6 6 1.41-1.41L10.83 12l4.58-4.59z"/></svg>
      </span>

      Back to {{ backTitle }}
    </button>

    <div class="mb-4">
      <patient-details :patient="patient" :editable="true"/>
    </div>
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
      editVisited: false,
      patient: {},
    }
  },
  computed: {
    backTitle() {
      switch (this.$route.name) {
        case 'ehr.patient.transfer.edit':
        case 'ehr.patient.edit':
          return 'patient overview';
        default:
          return 'patient list';
      }
    }
  },
  methods: {
    truncate(str, n) {
      return (str.length > n) ? str.substr(0, n - 1) + '...' : str
    },
    fetchPatient() {
      this.$api.getPatient({patientID: this.$route.params.id})
          .then(patient => this.patient = patient)
          .catch(error => this.$status.error(error))
    },
    back() {
      switch (this.$route.name) {
        case 'ehr.patient.transfer.edit':
        case 'ehr.patient.edit':
          this.$router.push({name: 'ehr.patient', params: {id: this.$route.params.id}})
          break;
        default:
          this.$router.push({name: 'ehr.patients'});
      }
    },

    // Fetch new patient data after the patient was updated
    updateAfterEdit() {
      if (this.editVisited && this.$route.name === 'ehr.patient.overview') {
        this.editVisited = false;
        this.fetchPatient();

        return;
      }

      if (this.$route.name === 'ehr.patient.edit') {
        this.editVisited = true;
      }
    }
  },
  created() {
    this.fetchPatient()
  },
  watch: {
    $route: 'updateAfterEdit'
  }
}
</script>
