<template>
  <div class="flex">
    <!-- Navigation -->
    <nav class="main-nav shadow-lg fixed flex flex-col justify-between top-0 h-full w-72 bg-white">
      <div class="px-3 py-8">
        <div class="px-3 mb-5 inline-flex items-center">
          <div class="w-10 mr-2">
            <svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" version="1.1" id="Layer_1" x="0px" y="0px" viewBox="0 0 272.1 119.8" style="enable-background:new 0 0 272.1 119.8;" xml:space="preserve">
              <path fill="#D15949" class="st0" d="M210.6,0h-82.1v24.8h29.8v94.9h24.8V24.8h9.3c1.6-7.4,5.4-14.2,11-19.4C205.6,3.4,208,1.5,210.6,0z"/>
              <path fill="#D15949" class="st0" d="M264.9,61.5c-4.8-5-14.1-10.5-28-16.5c-10.1-4.3-15.1-8.6-15.1-12.9c-0.1-2.3,0.9-4.5,2.6-6.1  c1.7-1.6,3.9-2.4,6.2-2.4l0,0c5,0,9.3,3.5,12.7,10.5l21.6-11.5c-7-15-17.8-22.4-32.3-22.5l0,0c-6.9-0.1-13.7,1.8-19.6,5.4  c-2,1.2-3.8,2.7-5.5,4.2c-4.3,4.1-7.4,9.3-8.9,15.1c-0.8,3.1-1.2,6.2-1.2,9.4c0,7.4,2,13.5,5.9,18.3c3.1,3.7,8.4,7.4,16,11.2  c2.7,1.3,5.4,2.6,8,3.9c2.6,1.2,5.3,2.5,7.9,3.8c6.8,4,10.2,8.3,10.2,13c0,8.1-4.4,12.1-13.2,12.1l0,0c-4.2,0.1-8.3-1.5-11.4-4.4  c-1.6-1.5-4-4.9-7.1-10.4l-22.3,10.3c7.9,18.5,21.3,27.8,40.3,27.8l0,0c11.1,0,20.5-3.5,28.4-10.5c7.9-6.9,12.3-17,12.1-27.4  C272.1,73.3,269.7,66.5,264.9,61.5z"/>
              <path fill="#D15949" class="st0" d="M128.5,30.8v42c0,17.1-6.7,25.6-20.2,25.6c-13.1,0-19.8-8-20.2-24.1V46.4c0-9.5-1-16.6-3-21.4  c-2.2-5.6-5.8-10.6-10.5-14.5C66.6,3.5,56.4,0,44,0C31.7,0,21.5,3.5,13.3,10.5C8.8,14.4,5.2,19.4,3.1,25C0.9,31.9-0.2,39.1,0,46.4  v73.3c0.6-0.7,1.3-1.4,2-2.1c6-5.7,13.7-9.4,21.9-10.6v-60c0-17.1,6.7-25.6,20.2-25.6c13.1,0,19.8,8,20.2,24.1v27.9  c0,9.5,1,16.6,3,21.4c2.2,5.6,5.8,10.6,10.4,14.5c8,7,18.2,10.5,30.6,10.5c12.4,0,22.6-3.5,30.7-10.5c4.6-3.9,8.1-8.9,10.3-14.5  c2.2-6.9,3.2-14.1,3.1-21.4V30.8H128.5z"/>
            </svg>
          </div>

          <div class="main-title text-xl">
            <small>/</small> <span class="font-semibold">Demo EHR</span>
          </div>
        </div>

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
            <inbox-badge v-if="inboxInfo && inboxInfo.messageCount > 0" :inbox="inboxInfo"/>
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
      <div class="w-full m-auto max-w-screen-2xl">
        <!-- Main content -->
        <div>
          <status-reporter type="error" :message="errorMessage.value"></status-reporter>
          <status-reporter type="info" :message="statusMessage.value"></status-reporter>
          <router-view @statusUpdate="updateStatus"></router-view>
        </div>

        <status-bar :statusMessage="eventMessage"></status-bar>
      </div>
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
  @apply rounded mb-2 py-3 px-3 inline-flex items-center text-gray-500 text-sm;
}

.menu-link:hover {
  @apply text-gray-800 bg-gray-100;
}

.menu-link-active, .menu-link-active:hover {
  @apply bg-gray-200 text-black;
}

.main-title {
  color: #D15949;
}
</style>
