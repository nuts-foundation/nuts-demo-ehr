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
            url: o => `${o.url}/session`,
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
              // Send the session token in the request which contains the current customerID
              'Authorization': `Bearer ${localStorage.getItem("session")}`
            },
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
      };
      let irmaPopup = irma.newPopup(options);
      irmaPopup.start()
          .then(responseData => {
            localStorage.setItem("session", responseData.token)
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
