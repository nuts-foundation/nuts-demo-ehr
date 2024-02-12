<template>
  <div class="flex justify-center">

    <div class="mt-12 bg-white border rounded-md max-w-7xl p-8 flex flex-col">
      <h1 class="text-3xl py-2">WebAuthn login</h1>
      <form class="my-4 flex justify-center" @submit.stop.prevent="login">
        <div class="space-y-4">

          <div class="text-sm font-medium text-gray-700">Organization: {{ customer.name }}</div>

          <div>
            <p>Please follow the instructions on screen.</p>
          </div>

          <p v-if="!!loginError" class="p-2 text-center bg-red-100 rounded-md">{{ loginError }}</p>
        </div>
      </form>
    </div>
  </div>
</template>

<script>
export default {
  props: ['redirectPath'],
  data() {
    return {
      loginError: "",
      credentials: {
        customerID: null
      },
      customer: {
        name: "",
      },
    }
  },
  watch: {
    // Remove error when typing
    // 'credentials.password'() {
    //   this.loginError = ""
    // },
  },
  methods: {
    redirectAfterLogin() {
      if (this.redirectPath) {
        return this.$router.push(this.redirectPath)
      }
      this.$router.push("/ehr/")
    },
    arrayBufferToBase64(buffer) {
      // Create a typed array from the buffer
      let bytes = new Uint8Array(buffer);
      // Convert the bytes to a string of characters
      let binary = "";
      for (let i = 0; i < bytes.length; i++) {
        binary += String.fromCharCode(bytes[i]);
      }
      // Encode the string using base64
      return window.btoa(binary);
    },
  },
  mounted() {
    // // Check if session still valid, if so just redirect to application
    // Currently disabled, make conditional on existence of cookie / JWT
    // this.$api.get('web/private')
    //     .then(() => this.redirectAfterLogin())
    //     .catch(() => {
    //       // session is invalid, need to authenticate
    //     }))
    if ('customer' in this.$route.params) {
      this.customer = JSON.parse(this.$route.params.customer)
      this.credentials.customerID = this.customer.id
    } else {
      // Missing required params, redirect to landing page
      console.log("missing customer in params. Back to login page.")
      this.$router.push("/")
    }
    // check for webauthn cookie ('webauthn')
    // if not present, show signup pop-up
    // if present, show login pop-up
    // get cookie named 'credential'
    let credentialString = localStorage.getItem('credential')
    if (credentialString === null) {
      // create new credential with "random" id
      const userID = "UZSL85T9AFC";

      // show signup pop-up
      // for now/demo we just show the WebAuthn dialog
      let publicKeyCredentialCreationOptions = {
        challenge: new Uint8Array([ // must be a cryptographically random number sent from a server
          0x8c, 0x0a, 0x0a, 0x0a, 0x0a, 0x0a, 0x8c, 0x0a, 0x0a, 0x0a, 0x0a, 0x0a,
          0x8c, 0x0a, 0x0a, 0x0a, 0x0a, 0x0a
        ]).buffer,
        rp: {
          name: 'Nuts Demo EHR.' // name of the relying party
        },
        user: {
          id: Uint8Array.from(userID, c => c.charCodeAt(0)),
          name: 'demo',
          displayName: "You",
        },
        pubKeyCredParams: [{
            type: 'public-key', // type of credential
            alg: -7 // ES256
          },
          {
            type: 'public-key', // type of credential
            alg: -257 // RS256
          }
        ],
        authenticatorSelection: {
          authenticatorAttachment: 'platform', // platform, cross-platform, or null
          requireResidentKey: false, // require that the authenticator have a resident key
          userVerification: 'preferred' // required, preferred, or discouraged
        },
        timeout: 60000, // time in milliseconds
        attestation: 'direct', // none, indirect, or direct
      }
      navigator.credentials.create({ publicKey: publicKeyCredentialCreationOptions })
        .then((newCredentialInfo) => {
          // send attestation response and user's email to server
          console.log(newCredentialInfo)
          // public key
          const derEncodedKey = newCredentialInfo.response.getPublicKey()
          const base64EncodedKey = this.arrayBufferToBase64(derEncodedKey)
          const credential = {
            userId: userID,
            credentialId: this.arrayBufferToBase64(newCredentialInfo.rawId),
            publicKey: base64EncodedKey,
          }
          console.log('storing credential', credential)
          localStorage.setItem('credential', JSON.stringify(credential))
          // submit to server
          this.$api.authenticateWithWebAuthn({body: {customerID: this.credentials.customerID, credentialID: credential.credentialId, userID: credential.userId, publicKey: credential.publicKey}})
              .then(responseData => {
                localStorage.setItem("session", responseData.token)
                console.log("WebAuthn authentication successful")
                this.redirectAfterLogin()
              })
              .catch(response => {
                console.log("WebAuthn authentication failed: " + response)
              })
        })
        .catch((err) => {
          console.log(err)
        })
    } else {
      let credential = JSON.parse(credentialString)
      let binaryString = atob(credential.credentialId)
      console.log(credential)
      const challenge = "123456"
      // do login
      const publicKeyCredentialRequestOptions = {
        challenge: Uint8Array.from(challenge, c => c.charCodeAt(0)),
        allowCredentials: [{
          id: Uint8Array.from(
              binaryString, c => c.charCodeAt(0)),
          type: 'public-key',
          transports: ['internal'],
        }],
        timeout: 60000,
      }
      const assertion = navigator.credentials.get({ publicKey: publicKeyCredentialRequestOptions })
        .then((assertion) => {
          // send assertion response to server for verification
          console.log(assertion)
          this.$api.authenticateWithWebAuthn({body: {customerID: this.credentials.customerID, credentialID: credential.credentialId, userID: credential.userId, publicKey: credential.publicKey}})
              .then(responseData => {
                localStorage.setItem("session", responseData.token)
                console.log("WebAuthn authentication successful")
                this.redirectAfterLogin()
              })
              .catch(response => {
                console.log("WebAuthn authentication failed: " + response)
              })
        })
        .catch((err) => {
          console.log(err)
        })
    }
  },
}
</script>
