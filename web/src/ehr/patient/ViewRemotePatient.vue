<template>
  <div class="px-12 py-8">
    <h1 class="mt-12">View remote patient file</h1>
    <p class="my-6">
      Here you can view the medical records of a patient at a remote care organization.
    </p>
    <section v-if="!patientFile" class="mt-8 bg-white p-5 shadow-lg rounded-lg">
      <div>
        <label for="organizationSelect">
          <div>Organization</div>
          <small class="font-light">Note that the organization must have authorized your care organization to access the
            patient's records.</small>
        </label>
        <select id="organizationSelect" v-model="chosenOrganization">
          <option :value="organization.did" v-for="organization in organizations" :key="organization.did">
            {{ organization.name + " (" + organization.discoveryServices.join(', ') + ")" }}
          </option>
        </select>
      </div>
      <div class="my-3">
        <label for="patientSSN">
          <div>Patient SSN</div>
          <small class="font-light">Enter the patient's social security number, which is used to look up the patient's
            records</small>
        </label>
        <input type="text" v-model="chosenPatientSSN" id="patientSSN">
      </div>
      <div class="my-3">
        <label for="requestedScope">
          <div>OAuth2 scope</div>
          <small class="font-light">Specifies the scope of the requested access</small>
        </label>
        <input type="text" v-model="requestedScope" id="requestedScope">
      </div>
      <div class="my-3">
        <button v-on:click="viewPatient" class="btn btn-primary">View patient</button>
      </div>
    </section>
    <section v-if="patientFile" class="mt-8 bg-white p-5 shadow-lg rounded-lg">
      <PatientDetails :patient="patientFile.patient"/>
      <div class="py-3">
        <FHIRObservations :observations="patientFile.observations"/>
      </div>
      <div class="my-3">
        <button v-on:click="this.patient = null" class="btn btn-primary">Back</button>
      </div>
    </section>
  </div>
</template>
<script>

import PatientDetails from "./PatientDetails.vue";
import FHIRObservations from "./FHIRObservations.vue";

export default {
  components: {FHIRObservations, PatientDetails},
  data() {
    return {
      loading: false,
      formErrors: [],
      patientFile: null,
      chosenPatientSSN: '1234567890',
      chosenOrganization: null,
      requestedScope: "homemonitoring",
      organizations: [],
    }
  },
  methods: {
    viewPatient() {
      this.patientFile = null
      this.$api.getRemotePatient({
        remotePartyDID: this.chosenOrganization,
        patientSSN: this.chosenPatientSSN,
        scope: this.requestedScope
      })
          .then((result) => {
            this.patientFile = result.data
      })
      .catch(error => this.$status.error(error))
    },
  },
  mounted() {
    this.$api.searchOrganizations(null, {"issuer": "*"})
        .then((results) => this.organizations = Object.values(results.data))
        .catch(error => this.$status.error(error))
  },
}
</script>
