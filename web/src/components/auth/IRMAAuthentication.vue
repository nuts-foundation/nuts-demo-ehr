<template>
  <div></div>
</template>
<script>
import irma from "@privacybydesign/irma-frontend";

export default {
  data() {
    return {
      customer: null,
    }
  },
  mounted() {
    if ('customer' in this.$route.params) {
      this.customer = JSON.parse(this.$route.params.customer)
      this.perform()
    } else {
      // Missing required params, redirect to landing page
      this.$router.push("/")
    }
  },
  methods: {
    redirectAfterLogin() {
      console.log('logged in, redirecting!')
      this.$router.push("/ehr/")
    },
    perform() {
      let options = {
        // Developer options
        debugging: true,

        // Front-end options
        language: 'en',
        translations: {
          header: "Authenticate for " + this.customer.name
        },

        // Back-end options
        session: {
          // Point to demo-ehr backend which forwards requests
          url: '/web/auth/irma',

          // Define your disclosure request:
          start: {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json'
            },
            body: JSON.stringify({customerID: this.customer.id})
          },
          mapping: {
            sessionPtr: r => r.sessionPtr.clientPtr,
            sessionToken: r => r.sessionID
          }
        }
      };
      let irmaPopup = irma.newPopup(options);
      irmaPopup.start()
          .then(result => {
            console.log("IRMA authentication successful")
            this.redirectAfterLogin()
          })
          .catch(error => {
            if (error === 'Aborted') {
              console.log('IRMA authentication aborted')
              this.$router.push("/")
              return;
            }
            console.error("error", error);
          })
          .finally(() => irmaPopup = irma.newPopup(options))
    }
  }
}
</script>