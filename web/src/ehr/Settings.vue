<template>
  <h1 class="page-title">Nuts Presence</h1>
  <p>Your organization registration on the Nuts network.</p>
  <div class=" mt-8">
    <div class="space-y-4 w-full" v-if="customer && customer.did">
      <form class="space-y-3">
        <div>
          <label for="did">DID</label>
          <input type="text" disabled v-model="customer.did" id="did">
        </div>
        <div>
          <label for="name">Name of the organization</label>
          <input type="text" disabled v-model="customer.name" id="name">
        </div>
        <div>
          <label for="city">City of registration</label>
          <input type="text" disabled v-model="customer.city" id="city">
        </div>
      </form>
    </div>
    <div class="space-y-4 w-full" v-if="!customer || customer.did == ''">
      <h2 class="page-subtitle">Your organization is not registered.</h2>
    </div>
  </div>
</template>
<script>
export default {
  data() {
    return {
      customer: null,
    }
  },
  created() {
    this.fetchData()
  },
  emits: ['statusUpdate'],
  methods: {
    fetchData() {
      this.$api.getCustomer()
          .then(responseData => this.customer = responseData)
          .catch(error => this.$status.report(error))
    }
  }
}
</script>
