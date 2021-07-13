<template>
  <div>
    <transfer-form v-if="transfer" :patient="transfer.patient" :transfer="transfer"
                   @input="(updatedTransfer) => {this.transfer = updatedTransfer}"/>

    <div class="mt-4">
      <button @click="updateTransfer" class="btn btn-primary">Update</button>
    </div>

    <table class="mt-4 min-w-full divide-y divide-gray-200" v-if="transfer">
      <thead class="bg-gray-50">
      <tr>
        <th>Organization</th>
        <th>Date</th>
        <th>Status</th>
      </tr>
      </thead>
      <tbody>
      <tr v-for="negotiation in negotiations">
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
            {{ slotProps.item.name }}
          </auto-complete>
        </td>
        <td colspan="2" v-if="!!requestedOrganization">
          {{ requestedOrganization.name }}
        </td>
        <td v-if="!!requestedOrganization">
          <button class="btn btn-primary" @click="startNegotiation">Request</button>
          <button class="btn" @click="cancelOrganization">Cancel</button>
        </td>
      </tr>
      </tbody>
    </table>

    <table class="min-w-full divide-y divide-gray-200 mt-4" v-if="transfer">
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
import TransferForm from "./TransferForm.vue"
import AutoComplete from "../../../components/Autocomplete.vue"

export default {
  components: {TransferForm, AutoComplete},
  data() {
    return {
      transfer: null,
      negotiations: [],
      messages: [
        {title: "Aanmeldbericht", contents: "Some content"},
        {title: "Overdrachtsbericht", contents: "Some content 2"},
      ],
      organizations: [
        {did: "1234", name: "HengeZorg", zipcode: "7552AB", starred: true},
        {did: "5678", name: "Zorgcentrum Enschede", zipcode: "7552CC", starred: false},
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
    startNegotiation() {
      this.$api.startTransferNegotiation({transferID: this.transfer.id, organizationDID: this.requestedOrganization.did})
    },
    updateTransfer() {

    },
    fetchTransfer(id) {
      this.$api.getTransfer({transferID: id})
          .then(transfer => this.transfer = transfer)
          .then(() => this.$api.listTransferNegotiations({transferID: id}))
          .then(negotiations => this.negotiations = negotiations)
          .catch(error => this.$errors.report(error))
    }
  },
  mounted() {
    this.fetchTransfer(this.$route.params.transferID)
  },

}
</script>
