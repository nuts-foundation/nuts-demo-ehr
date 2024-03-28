<template>
  <div class="px-12 py-8">
    <h1 class="mt-12">Nuts Presence</h1>

    <p>Your organization registration on the Nuts network.</p>

    <div class="mt-8 bg-white p-5 shadow-lg rounded-lg">
      <div class="space-y-4 w-full" v-if="customer && customer.did">
        <form class="space-y-3">
          <div>
            <label for="did">DID</label>
            <p id="did">{{customer.did}}</p>
          </div>
          <div>
            <label for="name">Name of the organization</label>
            <p id="name">{{customer.name}}</p>
          </div>
          <div>
            <label for="city">City of registration</label>
            <p id="city">{{customer.city}}</p>
          </div>
        </form>
      </div>
      <div class="space-y-4 w-full" v-if="!customer || customer.did == ''">
        <h2 class="page-subtitle">Your organization is not registered.</h2>
      </div>
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
  methods: {
    fetchData() {
      this.$api.getCustomer()
          .then(result => this.customer = result.data)
          .catch(error => this.$status.report(error))
    }
  }
}
</script>
