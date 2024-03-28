<template>
  <div>
    <h1>Edit Transfer</h1>
    <transfer-form
        v-if="transfer"
        :transfer="transfer"
        mode="edit"
        @input="(updatedTransfer) => {this.transfer = updatedTransfer}"/>

    <div class="bg-white p-5 shadow-sm rounded-lg mt-6">
      <table class="min-w-full divide-y divide-gray-200" v-if="transfer">
        <thead>
        <tr>
          <th>Organization</th>
          <th>Date</th>
          <th colspan="2">Status</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="negotiation in negotiations">
          <td v-if="negotiation.organization">{{ negotiation.organization.name }}</td>
          <td v-else>{{ negotiation.organizationDID }}</td>
          <td>{{ negotiation.transferDate }}</td>
          <td>
            <transfer-status :status="{status: negotiation.status}"/>
          </td>
          <td class="space-x-2">
            <span v-if="negotiation.status === 'accepted' && negotiation.status !== 'completed'"
                  @click="assignNegotiation(negotiation)" class="hover:underline cursor-pointer"
                  :class="{'btn-loading': state === 'assigning'}">assign</span>
            <span v-if="negotiation.status !== 'cancelled' && negotiation.status !== 'completed'"
                  @click="cancelNegotiation(negotiation)" class="hover:underline cursor-pointer"
                  :class="{'btn-loading': state === 'cancelling'}">cancel</span>
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
          <td v-if="!!requestedOrganization" class="space-x-2">
            <button class="btn btn-sm btn-primary" @click="assignOrganization" id="transfer-assign-button"
                    :class="{'btn-loading': state === 'assigning'}"><span>Assign<span style="font-family: monospace;" v-if="state === 'assigning'"> {{'.'.repeat(waitCount) + '&nbsp;'.repeat(3-waitCount)}}</span></span></button>
<!-- Not supported for now -->
<!--            <button class="btn btn-sm btn-primary" @click="startNegotiation"-->
<!--                    :class="{'btn-loading': state === 'requesting'}">Request-->
<!--            </button>-->
            <button class="btn btn-sm btn-secondary" @click="cancelOrganization">Cancel</button>
          </td>
        </tr>
        <tr v-if="showRequestNewOrganization()">
          <td colspan="3">
            <p>Note: only care organizations that accept patient transfers over the Nuts Network can be selected.</p>
          </td>
        </tr>
        </tbody>
      </table>
    </div>

    <div class="mt-4 space-x-2">
      <button v-if="showUpdateButton"
              @click="updateTransfer" class="btn btn-primary">
        Update
      </button>

      <button v-if="transfer && transfer.status !== 'cancelled' && transfer.status !== 'completed'"
              @click="cancelTransfer" class="btn"
              :class="{'btn-secondary': showUpdateButton, 'btn-primary': !showUpdateButton}"
      >
        Cancel transfer
      </button>

      <button @click="$router.push({name: 'ehr.patient', params: {id: $route.params.id } })"
              class="btn btn-secondary">
        Back
      </button>
    </div>

    <table v-if="transfer && transfer.messages && transfer.messages.length > 0"
           class="min-w-full divide-y divide-gray-200 mt-6">
      <thead>
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
import TransferForm from "./TransferFields.vue"
import AutoComplete from "../../components/Autocomplete.vue"
import TransferStatus from "../../components/TransferStatus.vue"

export default {
  components: {TransferForm, AutoComplete, TransferStatus},
  data() {
    return {
      state: 'init',
      transfer: null,
      negotiations: [],
      waitCount: 1,
      messages: [
        {title: "Aanmeldbericht", contents: "Some content"},
        {title: "Overdrachtsbericht", contents: "Some content 2"},
      ],
      organizations: [],
      requestedOrganization: null,
    }
  },
  computed: {
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
    chooseOrganization(organization) {
      this.requestedOrganization = organization
    },
    cancelOrganization() {
      this.requestedOrganization = null
    },
    searchOrganizations(query) {
      this.$api.searchOrganizations({query: query, discoveryServiceType: "eoverdracht_dev3", didServiceType: "eOverdracht-receiver"})
          .then((result) => {
            // Only show organizations that we aren't already negotiating with
            this.organizations = result.data.filter(i => this.negotiations.filter(n => i.did === n.organizationDID).length === 0)
          })
          .catch(error => this.$status.error(error))
    },
    assignOrganization() {
      this.state = 'assigning';
      this.waitCount = 1;

      let timer = setInterval(() => {
        if (this.waitCount === 3) {
          this.waitCount = 1;
          return;
        }

        this.waitCount++
      }, 1000);

      const negotiation = {
        transferID: this.transfer.id,
      };

      this.$api.assignTransferDirect(negotiation, {
        organizationDID: this.requestedOrganization.did
      })
          .then(() => {
            this.requestedOrganization = null
            this.$store.commit("statusUpdate", "Patient transfer assigned")
            this.fetchTransfer(this.transfer.id)
          })
          .catch(error => this.$status.error(error))
          .finally(() => {
            clearInterval(timer)
            this.state = 'done'
          })
    },
    startNegotiation() {
      this.state = 'requesting';

      const negotiation = {
        transferID: this.transfer.id,
      };

      this.$api.startTransferNegotiation(negotiation, {
        organizationDID: this.requestedOrganization.did
      })
          .then(() => {
            this.requestedOrganization = null
            this.fetchTransfer(this.transfer.id)
          })
          .catch(error => this.$status.error(error))
          .finally(() => this.state = 'done')
    },
    assignNegotiation(negotiation) {
      this.state = 'assigning'
      this.$api.updateTransferNegotiationStatus(
          {transferID: negotiation.transferID, negotiationID: negotiation.id},
          {status: 'in-progress'}
      )
          .then(() => this.fetchTransferNegotiations(this.transfer.id))
          .catch(error => this.$status.error(error))
          .finally(() => this.state = 'done')
    },
    cancelNegotiation(negotiation) {
      this.state = 'cancelling';

      this.$api.updateTransferNegotiationStatus(
          {transferID: negotiation.transferID, negotiationID: negotiation.id},
          {status: 'cancelled'}
      )
          .then(() => this.fetchTransferNegotiations(this.transfer.id))
          .catch(error => this.$status.error(error))
          .finally(() => this.state = 'done')
    },
    updateNegotiation(negotiation) {
    },
    cancelTransfer() {
      const cancelRequest = {
        transferID: this.transfer.id,
      }
      this.$api.cancelTransfer(cancelRequest)
          .then(result => {
            this.transfer = result.data
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
      }
      this.$api.updateTransfer(updateRequest, {
        description: this.transfer.description,
        transferDate: this.transfer.transferDate,
      })
          .then(result => {
            this.transfer = result.data
            this.$status.status("Transfer updated")
          })
          .catch(error => this.$status.error(error))

    },
    fetchTransfer(id) {
      this.$api.getTransfer({transferID: id})
          .then(result => this.transfer = result.data)
          .then(() => this.fetchTransferNegotiations(id))
          .catch(error => this.$status.error(error))
    },
    fetchTransferNegotiations(transferID) {
      return this.$api.listTransferNegotiations({transferID: transferID})
          .then(result => this.negotiations = result.data)
          .catch(error => this.$status.error(error))
    }
  },
  mounted() {
    this.fetchTransfer(this.$route.params.transferID)
  },
}
</script>
