<template>
  <div class="px-12 py-8">
    <h1 class="mt-12">JWT signing demo</h1>

    <p>Sign content with your WebAuthn device</p>

    <div class="mt-8 bg-white p-5 shadow-lg rounded-lg">
      <div class="space-y-4 w-full">
        <form class="space-y-3">
          <div>
            <label for="content">Content</label>
            <textarea id="content" type="text" v-model="content" placeholder="enter credentialSubject"/>
          </div>
          <div class="inline-flex mt-8">
            <button
                @click="sign(content)"
                class="float-right inline-flex items-center bg-nuts text-white rounded-lg justify-center shadow-md p-2"
            >
              Sign
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>
<script>
import base64 from 'base-64'

export default {
  data() {
    return {
      content: "{\"test\":true}"
    }
  },
  created() {
    //this.fetchData()
  },
  methods: {
    sign(content) {
      // create the header
      const header = {
        kid: "did:jwk:eyJjcnYiOiJQLTI1NiIsImt0eSI6IkVDIiwieCI6InFOT19sdFlOYTlpNXR0R1FKN1phbkNtVmVzOW5rRXlUUUs5bmZMYWFvNHMiLCJ5IjoiX2VDWWtkUUxBRXFrX3VJUWc2bGItdHFZM0tfRzE5UXZHeF91NnVjSTluVSJ9#0",
        alg: "ES256",
        typ: "JWT"
      }
      const payload = {
        iss: "did:jwk:eyJjcnYiOiJQLTI1NiIsImt0eSI6IkVDIiwieCI6InFOT19sdFlOYTlpNXR0R1FKN1phbkNtVmVzOW5rRXlUUUs5bmZMYWFvNHMiLCJ5IjoiX2VDWWtkUUxBRXFrX3VJUWc2bGItdHFZM0tfRzE5UXZHeF91NnVjSTluVSJ9",
        vc: {
          '@context': "https://www.w3.org/2018/credentials/v1",
          id: "did:jwk:eyJjcnYiOiJQLTI1NiIsImt0eSI6IkVDIiwieCI6InFOT19sdFlOYTlpNXR0R1FKN1phbkNtVmVzOW5rRXlUUUs5bmZMYWFvNHMiLCJ5IjoiX2VDWWtkUUxBRXFrX3VJUWc2bGItdHFZM0tfRzE5UXZHeF91NnVjSTluVSJ9#1",
          issuanceDate: "2023-11-06T15:00:00Z",
          type: "VerifiableCredential",
          credentialSubject: JSON.parse(content)
        },
      }
      console.log(header)
      console.log(payload)
      const challenge = base64.encode(JSON.stringify(header)).replace(/=+$/, '') + "." + base64.encode(JSON.stringify(payload)).replace(/=+$/, '')
      console.log(challenge)

      let credentialString = localStorage.getItem('credential')
      let credential = JSON.parse(credentialString)
      console.log(credential.credentialId)
      let binaryString = atob(credential.credentialId)
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
            const sigbytes = new Uint8Array(assertion.response.signature)
            const signature = base64.encode(sigbytes)
            console.log(signature)
          })
          .catch((err) => {
            console.log(err)
          })


    }
  }
}
</script>
