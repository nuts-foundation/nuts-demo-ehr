<template>
  <div>
    <p v-if="apiError" class="p-3 bg-red-100 rounded-md">Error: {{ apiError }}</p>
    <patient-details :patient="patient"/>

    <table class="mt-4 min-w-full divide-y divide-gray-200">
      <thead class="bg-gray-50">
      <tr>
        <th>Description</th>
      </tr>
      </thead>
      <tbody>
      <tr>
        <td><textarea v-model="description" class="border min-w-full h-32"></textarea></td>
      </tr>
      </tbody>
    </table>

    <table class="mt-4 min-w-full divide-y divide-gray-200">
      <thead class="bg-gray-50">
      <tr>
        <th>Organization</th>
        <th>Date</th>
        <th>Status</th>
      </tr>
      </thead>
      <tbody>
      <tr v-for="transfer in transfers">
        <td>{{ transfer.organization.name }}</td>
        <td>{{ transfer.date }}</td>
        <td>{{ transfer.status }}</td>
      </tr>
      <tr>
        <td colspan="3" v-if="requestedOrganization === null">
          <Autocomplete @input="searchOrganization" :results="organizations" @onSelect="chooseOrganization"
                        :input-class="['border']"
                        :results-container-class="['border']"
                        :results-item-class="['cursor-pointer', 'p-2', 'hover:bg-gray-100']"
          ></Autocomplete>
        </td>
        <td colspan="2" v-if="!!requestedOrganization">
          {{requestedOrganization.name}}
        </td>
        <td v-if="!!requestedOrganization">
          <button class="btn btn-primary">Request</button>
          <button class="btn" @click="cancelOrganization">Cancel</button>
        </td>
      </tr>
      </tbody>
    </table>

    <table class="min-w-full divide-y divide-gray-200 mt-4">
      <thead class="bg-gray-50">
      <tr>
        <th>Messages</th>
      </tr>
      </thead>
      <tbody>
      <tr v-for="message in messages">
        <td>{{ message.title }}</td>
      </tr>
      </tbody>
    </table>
  </div>
</template>
<script>
import Autocomplete from 'vue3-autocomplete'
import PatientDetails from "../PatientDetails.vue";

export default {
  components: {PatientDetails, Autocomplete},
  data() {
    return {
      apiError: null,
      patient: null,
      description: "Meneer heeft wondzorg nodig aan rechterbeen. 3 maal daags verband wisselen.",
      transfers: [
        {date: "2021-06-22", status: "in afwachting", organization: {name: "De Regenboog"}},
        {date: "2021-06-23", status: "geaccepteerd", organization: {name: "Avondrust"}},
      ],
      messages: [
        {title: "Aanmeldbericht", contents: "Some content"},
        {title: "Overdrachtsbericht", contents: "Some content 2"},
      ],
      organizations: [],
      requestedOrganization: null,
    }
  },
  methods: {
    searchOrganization(query) {
      console.log('searching for ' + query)
      this.organizations = [
        {name: "HengeZorg", zipcode: "7552AB", starred: true},
        {name: "Zorgcentrum Enschede", zipcode: "7552CC", starred: false},
      ]
    },
    chooseOrganization(organization) {
      this.requestedOrganization = organization
    },
    cancelOrganization() {
      this.requestedOrganization = null
    }
  },
  mounted() {
    this.$api.get(`/web/private/patient/${this.$route.params.id}`)
        .then(patient => this.patient = patient)
        .catch(error => this.apiError = error)
  }
}
</script>
