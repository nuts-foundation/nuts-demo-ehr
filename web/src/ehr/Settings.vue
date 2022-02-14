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
            <p id="name">{{kvkDetails ? kvkDetails.legalEntity : customer.name}}</p>
          </div>
          <div>
            <label for="city">City of registration</label>
            <p id="city">{{kvkDetails ? kvkDetails.officeAddress : customer.city}}</p>
          </div>
        </form>
      </div>
      <div class="space-y-4 w-full" v-if="!customer || customer.did == ''">
        <h2 class="page-subtitle">Your organization is not registered.</h2>
      </div>
      <div class="mt-8 w-full" v-if="kvkDetails === null || !kvkDetails.valid">
        <h3 class="font-semibold">IRMA verification</h3>

        <button @click="verify" class="btn btn-primary mt-3">Verify</button>
      </div>
    </div>
  </div>
</template>
<script>
import irma from "@privacybydesign/irma-frontend";

export default {
  data() {
    return {
      customer: null,
      kvkDetails: null,
    }
  },
  created() {
    this.fetchData()
  },
  methods: {
    fetchData() {
      this.$api.getCustomer()
          .then(responseData => this.customer = responseData)
          .catch(error => this.$status.report(error))

      this.$api.getKVKDetails()
          .then(responseData => this.kvkDetails = responseData)
          .catch(error => this.$status.report(error))
    },
    verify() {
      const popup = irma.newPopup({
        debugging: true,
        language: 'en',
        session: {
          url: "/web/auth/irma/kvk",
          start: {
            url: o => `${o.url}/session`,
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
              // Send the session token in the request which contains the current customerID
              'Authorization': `Bearer ${localStorage.getItem("session")}`
            }
          },
          mapping: {
            sessionPtr: r => r.sessionPtr.clientPtr,
            sessionToken: r => r.sessionID
          },
          result: {
            url: (o, {sessionPtr, sessionToken}) => `${o.url}/session/${sessionToken}/result`,
            headers: {
              'Authorization': `Bearer ${localStorage.getItem("session")}`
            },
            parseResponse: r => r.json()
          }
        }
      });

      popup.start()
        .then(responseData => {
          console.log("TOKEN", responseData)
        });
    }
  }
}
</script>
