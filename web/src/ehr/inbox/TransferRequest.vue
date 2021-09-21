<template>
  <div v-if="transferRequest">
    <h1 class="page-title">Transfer Request</h1>

    <div>
      <div class="mt-4">
        <div class="bg-gray-50 font-bold">Requesting Care Organization</div>
        <div>
          {{ transferRequest.sender.name }}, {{ transferRequest.sender.city }}
        </div>
      </div>

      <div class="mt-4">
        <div class="bg-gray-50 font-bold">Status</div>
        <div>
          {{ transferRequest.status }}
        </div>
      </div>

      <div class="mt-4">
        <div class="bg-gray-50 font-bold">Transfer date</div>
        <div>
          {{ transferRequest.advanceNotice.transferDate }}
        </div>
      </div>

      <div class="mt-4">
        <div class="bg-gray-50 font-bold">Condition</div>
        <h1>Problems:</h1>
        <ul class="list-decimal pl-4">
          <li
              v-for="patientProblem in transferRequest.advanceNotice.carePlan.patientProblems"
              class="pl-4"
          >
            <p> {{ patientProblem.problem.name }} </p>
            <div>
              <h2>Interventions:</h2>
              <ul class="list-disc">
                <li v-for="intervention in patientProblem.interventions">
                  {{ intervention.comment }}
                </li>
              </ul>
            </div>
          </li>
        </ul>
      </div>
    </div>

    <div class="mt-4">
      <button class="btn btn-primary m-1" @click="complete" v-show="transferRequest.status == 'in-progress'">Complete
      </button>
      <button class="btn btn-primary m-1" @click="accept" v-show="transferRequest.status == 'requested'">Accept</button>
      <button class="btn btn-secondary m-1" @click="reject" v-show="transferRequest.status == 'requested'">Reject
      </button>
    </div>
  </div>
</template>
<script>
export default {
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
