<template>
  <div class="px-12 py-8">
    <button type="button" @click="() => this.$router.push({name: 'ehr.inbox'})" class="btn btn-link mb-12">
      <span class="w-6 mr-1">
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="#000000"><path d="M0 0h24v24H0V0z"
                                                                                         fill="none"/><path
            d="M15.41 7.41L14 6l-6 6 6 6 1.41-1.41L10.83 12l4.58-4.59z"/></svg>
      </span>

      Back to inbox
    </button>

    <div v-if="transferRequest">
      <div class="mb-6 flex items-center justify-between">
        <h1>Transfer Request</h1>
      </div>

      <div class="bg-white rounded-lg shadow-lg">
        <div class="p-5 bg-gray-50 border-b rounded-t-lg">
          <div class="flex justify-between items-center">
            <div>
              <h2 class="mb-1">Requesting Care Organization</h2>

              <div class="text-gray-700" id="requesting-care-organization-info">
                {{ transferRequest.sender.name }} in {{ transferRequest.sender.city }}
              </div>
            </div>

            <div id="transfer-request-status-info">
              <transfer-status :status="{status: transferRequest.status}"/>
            </div>
          </div>
        </div>

        <div class="p-5" v-if="transferRequest.status != 'completed'">
          <div v-if="transferRequest.nursingHandoff">
            <div>
              <label>Transfer date</label>
              <div id="transfer-request-date-info">
                {{ transferRequest.nursingHandoff.transferDate }}
              </div>
            </div>

            <div class="mt-4">
              <label>Problems</label>

              <ul>
                <li v-for="patientProblem in transferRequest.nursingHandoff.carePlan.patientProblems">
                  <h3 class="font-semibold text-sm">Problem</h3>

                  <p data-problem-detail="name"> {{ patientProblem.problem.name }} </p>

                  <div class="mt-2">
                    <h3 class="font-semibold text-sm">Interventions</h3>

                    <ul>
                      <li v-for="intervention in patientProblem.interventions">
                        - &nbsp;<span data-problem-detail="intervention">{{ intervention.comment }}</span>
                      </li>
                    </ul>
                  </div>
                </li>
              </ul>
            </div>
          </div>
          <div v-else-if="transferRequest.advanceNotice && transferRequest.status != 'completed'">
            <div>
              <label>Transfer date</label>
              <div id="transfer-request-date-info">
                {{ transferRequest.advanceNotice.transferDate }}
              </div>
            </div>

            <div class="mt-4">
              <label>Problems</label>

              <ul>
                <li v-for="patientProblem in transferRequest.advanceNotice.carePlan.patientProblems">
                  <h3 class="font-semibold text-sm">Problem</h3>

                  <p data-problem-detail="name"> {{ patientProblem.problem.name }} </p>

                  <div class="mt-2">
                    <h3 class="font-semibold text-sm">Interventions</h3>

                    <ul>
                      <li v-for="intervention in patientProblem.interventions">
                        - &nbsp;<span data-problem-detail="intervention">{{ intervention.comment }}</span>
                      </li>
                    </ul>
                  </div>
                </li>
              </ul>
            </div>
          </div>
        </div>
      </div>

      <div v-if="transferRequest.status != 'completed'">

        <h2 class="mt-10 mb-3" v-if="transferRequest">Patient</h2>

        <div class="bg-white rounded-lg shadow-lg p-5">
          <div v-if="!transferRequest.nursingHandoff">
            <label>Zipcode</label>

            <div id="patient-zipcode-label">{{ transferRequest.advanceNotice.patient.zipcode }}</div>
          </div>

          <div v-if="transferRequest.nursingHandoff">
            <patient-details :patient="transferRequest.nursingHandoff.patient"/>
          </div>
        </div>

        <div class="mt-10">
          <button class="btn btn-primary" @click="complete" :class="{'btn-loading': state === 'completing'}"
                  v-show="transferRequest.status === 'in-progress'">
            Finish Transfer
          </button>

          <button class="btn btn-primary m-1" @click="accept" :class="{'btn-loading': state === 'accepting'}"
                  v-show="transferRequest.status === 'requested'">Accept
          </button>
          <button class="btn btn-secondary m-1" @click="reject" v-show="transferRequest.status === 'requested'">Reject
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
<script>


import TransferStatus from "../../components/TransferStatus.vue";
import PatientDetails from "../patient/PatientDetails.vue";

export default {
  components: {
    PatientDetails,
    TransferStatus
  },
  data() {
    return {
      state: 'init',
      transferRequest: null,
      token: null,
    }
  },
  created() {
    // when successfully authenticated, fetchData is called
    this.authenticate()
  },
  methods: {
    fetchData(token) {
      this.$api.getTransferRequest({
        requestorDID: this.$route.params.requestorDID,
        fhirTaskID: this.$route.params.fhirTaskID,
        token: token
      })
          .then((result) => this.transferRequest = result.data)
          .catch(error => this.$status.error(error))
    },
    accept() {
      this.state = 'accepting';

      this.$api.changeTransferRequestState({
        requestorDID: this.$route.params.requestorDID,
        fhirTaskID: this.$route.params.fhirTaskID,
      }, {status: 'accepted'})
          .then(() => this.fetchData())
          .catch(error => this.$status.error(error))
          .finally(() => this.state = 'done')
    },
    complete() {
      this.state = 'completing';

      this.$api.changeTransferRequestState(
          {
            requestorDID: this.$route.params.requestorDID,
            fhirTaskID: this.$route.params.fhirTaskID
          },
          {status: 'completed'}
      )
          .then(() => this.fetchData())
          .catch(error => this.$status.error(error))
          .finally(() => this.state = 'done')
    },
    reject() {
      alert('Not implemented yet.')
    },
    authenticate() {
      // the other side
      const requestorDID = this.$route.params.requestorDID
      // set closingUrl to the current root domain appended with #/close
      const closingUrl = window.location.origin + "/#/close"
      this.$api.createAuthorizationRequest({verifier: requestorDID, scope: "eOverdracht-sender", redirect_uri: closingUrl})
          .then(result => {
            // show a popup using the redirect URL
            window.open(result.data.redirect_uri, "_blank", "width=400,height=600")
            console.log("session_id: " + result.data.session_id)

            // let's poll the server for the result
            // loop and call every 1 second
            const session_id = result.data.session_id
            let interval = setInterval(() => {
              this.$api.getOpenID4VPAuthenticationResult({token: session_id})
                  .then(result => {
                    // check if result.data.access_token is available
                    if (result.data.access_token) {
                      clearInterval(interval)
                      console.log("OpenID4VP authentication successful")
                      console.log("AccessToken: " + result.data.access_token)
                      this.fetchData(result.data.access_token)
                    }
                  })
                  .catch(error => {
                    console.log("OpenID4VP authentication failed: " + error)
                    //this.loginError = response
                  })
            }, 20)

            // todo: continue to fetch data using the access token
            //localStorage.setItem("session", responseData.token)
            //this.redirectAfterLogin()
          })
          .catch(error => {
            console.log("OpenID4VP authentication failed: " + error)
            this.loginError = error
          })
    },
  }
}
</script>
