<template>
  <div>
    <h1 class="page-title">Transfer Request</h1>

    <div class="mt-4" v-if="transferRequest">
      <div class="bg-gray-50 font-bold">Requesting Care Organization</div>
      <div>
        {{ transferRequest.sender.name }}, {{ transferRequest.sender.city }}
      </div>
    </div>

    <div class="mt-4" v-if="transferRequest">
      <div class="bg-gray-50 font-bold">Transfer date</div>
      <div>
        {{ transferRequest.transferDate }}
      </div>
    </div>

    <div class="mt-4" v-if="transferRequest">
      <div class="bg-gray-50 font-bold">Condition</div>
      <div>
        {{ transferRequest.description }}
      </div>
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
      this.$api.getTransferRequest({requestorDID: this.$route.params.requestorDID, fhirTaskID: this.$route.params.fhirTaskID})
          .then((transferRequest) => this.transferRequest = transferRequest)
          .catch(error => this.$status.error(error))
    }
  }
}
</script>
