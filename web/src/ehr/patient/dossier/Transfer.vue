<template>
  <div>
    <p v-if="apiError" class="p-3 bg-red-100 rounded-md">Error: {{ apiError }}</p>
    <patient-details :patient="patient"/>

    <div class="mt-4">
      <div class="bg-gray-50 font-bold">Description</div>
      <div>
        <textarea v-model="!!transfer.description" class="border min-w-full h-32"></textarea>
      </div>
    </div>
    <div class="mt-4">
      <div class="bg-gray-50 font-bold">Transfer date</div>
      <div>
        <td><input type="date" v-model="!!transfer.transferDate"></td>
      </div>
    </div>

    <div class="mt-4">
      <button @click="createTransfer" v-if="transfer === null" class="btn btn-primary">Create Transfer</button>
      <button @click="updateTransfer" v-if="transfer !== null" class="btn btn-primary">Update Transfer</button>
    </div>

    <table class="mt-4 min-w-full divide-y divide-gray-200" v-if="!!transfer">
      <thead class="bg-gray-50">
      <tr>
        <th>Organization</th>
        <th>Date</th>
        <th>Status</th>
      </tr>
      </thead>
      <tbody>
      <tr v-for="negotiation in transfer.negotiations">
        <td>{{ negotiation.organization.name }}</td>
        <td>{{ negotiation.date }}</td>
        <td>{{ negotiation.status }}</td>
      </tr>
      <tr>
        <td colspan="3" v-if="requestedOrganization === null">

          <auto-complete
              :items="organizations"
              v-model:selected="requestedOrganization"
              v-slot="slotProps"
          >
            {{slotProps.item.name}}
          </auto-complete>
        </td>
        <td colspan="2" v-if="!!requestedOrganization">
          {{ requestedOrganization.name }}
        </td>
        <td v-if="!!requestedOrganization">
          <button class="btn btn-primary">Request</button>
          <button class="btn" @click="cancelOrganization">Cancel</button>
        </td>
      </tr>
      </tbody>
    </table>

    <table class="min-w-full divide-y divide-gray-200 mt-4" v-if="!!transfer">
      <thead class="bg-gray-50">
      <tr>
        <th>Messages</th>
      </tr>
      </thead>
      <tbody>
      <tr v-for="message in transfer.messages">
        <td>{{ message.title }}</td>
      </tr>
      </tbody>
    </table>
  </div>
</template>
<script>
import PatientDetails from "../PatientDetails.vue";
import AutoComplete from "../../../components/Autocomplete.vue";

export default {
  components: {PatientDetails, AutoComplete},
  data() {
    return {
      apiError: null,
      patient: {},
      transfer: null,
      // description: "Meneer heeft wondzorg nodig aan rechterbeen. 3 maal daags verband wisselen.",
      // transfers: [
      //   {date: "2021-06-22", status: "in afwachting", organization: {name: "De Regenboog"}},
      //   {date: "2021-06-23", status: "geaccepteerd", organization: {name: "Avondrust"}},
      // ],
      messages: [
        {title: "Aanmeldbericht", contents: "Some content"},
        {title: "Overdrachtsbericht", contents: "Some content 2"},
      ],
      organizations: [
        {name: "HengeZorg", zipcode: "7552AB", starred: true},
        {name: "Zorgcentrum Enschede", zipcode: "7552CC", starred: false},
      ],
      requestedOrganization: null,
    }
  },
  methods: {
    chooseOrganization(organization) {
      this.requestedOrganization = organization
    },
    cancelOrganization() {
      this.requestedOrganization = null
    },
    createTransfer() {
      this.$api.createTransfer()
    },
    updateTransfer() {

    },
    fetchPatient(patientID) {
      this.$api.getPatient({patientID: patientID})
          .then(patient => this.patient = patient)
          .catch(error => this.apiError = error)
    }
  },
  mounted() {
    this.fetchPatient(this.$route.params.id)
  }
}
</script>
