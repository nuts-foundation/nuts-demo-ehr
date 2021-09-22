<template>
  <div class="px-12 py-8">
    <button type="button" @click="() => this.$router.push({name: 'ehr.inbox'})" class="btn btn-link mb-12">
      <span class="w-6 mr-1">
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="#000000"><path d="M0 0h24v24H0V0z"
                                                                                         fill="none"/><path
            d="M15.41 7.41L14 6l-6 6 6 6 1.41-1.41L10.83 12l4.58-4.59z"/></svg>
      </span>

      Back to inbox
    </button>

    <div v-if="transferRequest">
      <div class="mb-6 flex items-center justify-between">
        <h1>Transfer Request</h1>
      </div>

      <div class="bg-white rounded-lg shadow-lg">
        <div class="p-5 bg-gray-50 border-b rounded-t-lg">
          <h2>Requesting Care Organization</h2>

          <div class="text-gray-700">
            {{ transferRequest.sender.name }} in {{ transferRequest.sender.city }}
          </div>
        </div>

        <div class="p-5">
          <div>
            <label>Status</label>
            <div>
              <transfer-status :status="{status: transferRequest.status}"/>
            </div>
          </div>

          <div class="mt-4">
            <label>Transfer date</label>

            <div>
              {{ transferRequest.advanceNotice.transferDate }}
            </div>
          </div>

          <div class="mt-4">
            <label>Problems</label>

            <ul>
              <li v-for="patientProblem in transferRequest.advanceNotice.carePlan.patientProblems">
                <h3 class="font-semibold">Problem</h3>

                <p> {{ patientProblem.problem.name }} </p>

                <div class="mt-2">
                  <h3 class="font-semibold">Interventions</h3>

                  <ul>
                    <li v-for="intervention in patientProblem.interventions">
                      - &nbsp;{{ intervention.comment }}
                    </li>
                  </ul>
                </div>
              </li>
            </ul>
          </div>
        </div>
      </div>

      <div class="mt-4">
        <button class="btn btn-primary m-1" @click="complete" v-show="transferRequest.status == 'in-progress'">
          Complete
        </button>
        <button class="btn btn-primary m-1" @click="accept" v-show="transferRequest.status == 'requested'">Accept</button>
        <button class="btn btn-secondary m-1" @click="reject" v-show="transferRequest.status == 'requested'">Reject</button>
      </div>
    </div>
  </div>
</template>
<script>
import TransferStatus from "../../components/TransferStatus.vue";

export default {
  components: {
    TransferStatus
  },
  data() {
    return {
      transferRequest: null,
    }
  },
  created() {
    this.fetchData()
  },
  methods: {
    fetchData() {
      this.$api.getTransferRequest({
        requestorDID: this.$route.params.requestorDID,
        fhirTaskID: this.$route.params.fhirTaskID
      })
          .then((transferRequest) => this.transferRequest = transferRequest)
          .catch(error => this.$status.error(error))
    },
    accept() {
      this.$api.changeTransferRequestState({
        requestorDID: this.$route.params.requestorDID,
        fhirTaskID: this.$route.params.fhirTaskID,
        body: {status: 'accepted'}
      })
          .then(() => this.fetchData())
          .catch(error => this.$status.error(error))
    },
    complete() {
      this.$api.changeTransferRequestState({
        requestorDID: this.$route.params.requestorDID,
        fhirTaskID: this.$route.params.fhirTaskID,
        body: {status: 'completed'}
      })
          .then(() => this.fetchData())
          .catch(error => this.$status.error(error))
    },
    reject() {
      alert('Not implemented yet.')
    }
  }
}
</script>
