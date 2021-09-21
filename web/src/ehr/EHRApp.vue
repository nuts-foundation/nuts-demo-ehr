<template>
  <div class="flex">
    <!-- Navigation -->
    <nav class="fixed flex flex-col justify-between top-0 h-full w-72 bg-blue-800">
      <div class="px-3 py-4">
        <h2 class="px-3 text-white mb-4">Nuts EHR</h2>

        <div class="grid grid-cols-1">
          <router-link
              :to="{name: 'ehr.patients'}"
              class="menu-link"
              active-class="menu-link-active">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-3" fill="none" viewBox="0 0 24 24"
                 stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"/>
            </svg>
            Patients
          </router-link>
          <router-link
              :to="{name: 'ehr.inbox'}"
              class="menu-link"
              active-class="menu-link-active">
            <inbox-badge :inbox="inboxInfo"/>
            Inbox
          </router-link>
          <router-link
              :to="{name: 'ehr.settings'}"
              class="menu-link"
              active-class="menu-link-active">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-3" fill="none" viewBox="0 0 24 24"
                 stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4"/>
            </svg>
            Settings
          </router-link>
        </div>
      </div>

      <div class="px-3 py-4">
        <!--<div class="flex flex-row items-center">
          <div class="rounded-full mb-1 mx-3 bg-gray-200 h-10 w-10 text-sm flex items-center justify-center">
            Hi
          </div>-->

          <div class="grid grid-cols-1">
            <router-link
                to="/logout"
                active-class="menu-link-active"
                class="menu-link">
              Logout
            </router-link>
          </div>
        </div>
      <!--</div>-->
    </nav>

    <main class="ml-72 mb-14 w-full">
      <!-- Main content -->
      <div>
        <status-reporter type="error" :message="errorMessage.value"></status-reporter>
        <status-reporter type="info" :message="statusMessage.value"></status-reporter>
        <router-view @statusUpdate="updateStatus"></router-view>
      </div>

      <status-bar :statusMessage="eventMessage"></status-bar>
    </main>
  </div>
</template>

<script>
import StatusBar from "../components/StatusBar.vue"
import StatusReporter from "../components/StatusReporter.vue"
import InboxBadge from "../components/InboxBadge.vue"

export default {
  components: {StatusBar, StatusReporter, InboxBadge},
  data() {
    return {
      inboxInfo: null,
      eventMessage: '',
    }
  },
  methods: {
    updateStatus(status) {
      this.eventMessage = status
    },
    fetchData() {
      this.$api.getInboxInfo()
          .then((response) => this.inboxInfo = response)
          .catch(error => this.$status.error(error))
    },
  },
  created() {
    this.fetchData()
  },
  computed: {
    errorMessage() {
      return this.$status.errorMessage
    },
    statusMessage() {
      return this.$status.statusMessage
    }
  },
  beforeRouteUpdate(to, from, next) {
    this.$status.clearError()
    next()
  }
}
</script>

<style>
.menu-link {
  @apply rounded mb-2 py-2 px-3 inline-flex items-center text-blue-300 text-sm;
}

.menu-link:hover {
  @apply text-blue-900 bg-blue-500;
}

.menu-link-active, .menu-link-active:hover {
  @apply bg-blue-900 text-white;
}
</style>
