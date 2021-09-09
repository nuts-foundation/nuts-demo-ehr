<template>
  <div></div>
</template>
<script>
import irma from "@privacybydesign/irma-frontend";

export default {
  props: {
    headerMessage: {
      type: String,
      default: "Authenticate",
    }
  },
  mounted() {
    this.perform()
  },
  emits: ['success', 'aborted', 'error'],
  methods: {
    perform() {
      let options = {
        // Developer options
        debugging: true,

        // Front-end options
        language: 'en',
        translations: {
          header: this.headerMessage
        },

        // Back-end options
        session: {
          // Set base-url path to demo-ehr endpoint which forwards requests to IRMA server (embedded in the Nuts Node)
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
            console.log("IRMA authentication successful")
            this.$emit('success', responseData.token)
          })
          .catch(error => {
            if (error === 'Aborted') {
              console.log('IRMA authentication aborted')
              this.$emit('aborted', error)
              return;
            }
            console.error("error", error);
            this.$emit('error', error)
          })
          .finally(() => irmaPopup = irma.newPopup(options))
    }
  }
}
</script>
