<template>
  <div>
    <div class="mt-4" v-if="transfer">
      <div class="bg-gray-50 font-bold">State</div>
      <div>
        {{ transfer.status }}
      </div>
    </div>
    <transfer-form v-if="transfer" :transfer="transfer"
                   @input="(updatedTransfer) => {this.transfer = updatedTransfer}"/>

    <table class="mt-4 min-w-full divide-y divide-gray-200" v-if="transfer">
      <thead class="bg-gray-50">
      <tr>
        <th>Organizations</th>
        <th>Date</th>
        <th colspan="2">Status</th>
      </tr>
      </thead>
      <tbody>
      <tr v-for="negotiation in negotiations">
        <td v-if="negotiation.organization">{{ negotiation.organization.name }}</td>
        <td v-else>{{ negotiation.organizationDID }}</td>
        <td>{{ negotiation.transferDate }}</td>
        <td>{{ negotiation.status }}</td>
        <td class="space-x-2">
          <span v-if="negotiation.status != 'cancelled' && negotiation.status != 'completed'"
                @click="cancelNegotiation(negotiation)" class="hover:underline cursor-pointer">cancel</span>
          <!--          <span @click="updateNegotiation(negotiation)" class="hover:underline cursor-pointer">update</span>-->
        </td>
      </tr>
      <tr v-if="showRequestNewOrganization()">
        <td colspan="3" v-if="requestedOrganization === null">
          <auto-complete
              :items="organizations"
              @selected="chooseOrganization"
              @search="searchOrganizations"
              v-slot="slotProps"
          >
            {{ slotProps.item.name }}
          </auto-complete>
        </td>
        <td colspan="2" v-if="!!requestedOrganization">
          {{ requestedOrganization.name }}
        </td>
        <td v-if="!!requestedOrganization">
          <button class="btn btn-primary" @click="assignOrganization">Assign</button>
          <button class="btn" @click="startNegotiation">Request</button>
          <button class="btn" @click="cancelOrganization">Cancel</button>
        </td>
      </tr>
      <tr v-if="showRequestNewOrganization()">
        <td colspan="3">
          <p>Note: only care organizations that accept patient transfers over the Nuts Network can be selected.</p>
        </td>
      </tr>
      </tbody>
    </table>

    <div class="mt-4 space-x-2">
      <button v-if="showUpdateButton()" @click="updateTransfer" class="btn btn-primary">Update</button>
      <button v-if="transfer && transfer.status != 'cancelled'" @click="cancelTransfer" class="btn">Cancel transfer
      </button>
      <button @click="$router.push({name: 'ehr.patient', params: {id: $route.params.id } })"
              class="btn btn-secondary"
      >
        Back
      </button>
    </div>

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
      organizations: [],
      requestedOrganization: null,
    }
  },
  methods: {
    showRequestNewOrganization() {
      switch (this.transfer.status) {
        case 'cancelled':
        case 'completed':
        case 'assigned':
          return false
        default:
          return true
      }
    },
    showUpdateButton() {
      if (!this.transfer) {
        return false
      }
      switch (this.transfer.status) {
        case 'cancelled':
        case 'completed':
          return false
        default:
          return true
      }
    },
    chooseOrganization(organization) {
      this.requestedOrganization = organization
    },
    cancelOrganization() {
      this.requestedOrganization = null
    },
    searchOrganizations(query) {
      this.$api.searchOrganizations({query: query, didServiceType: "eOverdracht-receiver"})
          .then((organizations) => {
            // Only show organizations that we aren't already negotiating with
            this.organizations = organizations.filter(i => this.negotiations.filter(n => i.did === n.organizationDID).length === 0)
          })
          .catch(error => this.$status.error(error))
    },
    assignOrganization() {
      const negotiation = {
        transferID: this.transfer.id,
        body: {
          organizationDID: this.requestedOrganization.did
        }
      };
      this.$api.assignTransferDirect(negotiation)
          .then(() => {
            this.requestedOrganization = null
            this.fetchTransfer(this.transfer.id)
          })
    },
    startNegotiation() {
      const negotiation = {
        transferID: this.transfer.id,
        body: {
          organizationDID: this.requestedOrganization.did
        }
      };
      this.$api.startTransferNegotiation(negotiation)
          .then(() => {
            this.requestedOrganization = null
            this.fetchTransfer(this.transfer.id)
          })
    },
    cancelNegotiation(negotiation) {
      this.$api.updateTransferNegotiationStatus({
        transferID: negotiation.transferID,
        negotiationID: negotiation.id,
        body: {status: 'cancelled'}
      })
          .then(() => this.fetchTransferNegotiations(this.transfer.id))
    },
    updateNegotiation(negotiation) {
    },
    cancelTransfer() {
      const cancelRequest = {
        transferID: this.transfer.id,
      }
      this.$api.cancelTransfer(cancelRequest)
          .then(transfer => {
            this.transfer = transfer
            return this.fetchTransferNegotiations(this.transfer.id)
          })
          .then(() => {
            this.$status.status("Transfer cancelled")
          })
          .catch(error => this.$status.error(error))
    },
    updateTransfer() {
      const updateRequest = {
        transferID: this.transfer.id,
        body: {
          description: this.transfer.description,
          transferDate: this.transfer.transferDate,
        }
      };
      this.$api.updateTransfer(updateRequest)
          .then(transfer => {
            this.transfer = transfer
            this.$status.status("Transfer updated")
          })
          .catch(error => this.$status.error(error))

    },
    fetchTransfer(id) {
      this.$api.getTransfer({transferID: id})
          .then(transfer => this.transfer = transfer)
          .then(() => this.fetchTransferNegotiations(id))
          .catch(error => this.$status.error(error))
    },
    fetchTransferNegotiations(transferID) {
      return this.$api.listTransferNegotiations({transferID: transferID})
          .then(negotiations => this.negotiations = negotiations)
          .catch(error => this.$status.error(error))
    }
  },
  mounted() {
    this.fetchTransfer(this.$route.params.transferID)
  },

}
</script>
