<template>
  <div class="flex">
    <!-- Navigation -->
    <nav class="main-nav shadow-lg fixed flex flex-col justify-between top-0 h-full w-72 bg-white">
      <div>
        <div class="px-6 pt-8">
          <div class="flex items-center">
            <div class="w-10 mr-2">
              <svg id="Layer_1" x="0px" y="0px" viewBox="0 0 272.1 119.8" style="enable-background:new 0 0 272.1 119.8;"
                   xml:space="preserve">
              <path fill="#D15949" class="st0"
                    d="M210.6,0h-82.1v24.8h29.8v94.9h24.8V24.8h9.3c1.6-7.4,5.4-14.2,11-19.4C205.6,3.4,208,1.5,210.6,0z"/>
                <path fill="#D15949" class="st0"
                      d="M264.9,61.5c-4.8-5-14.1-10.5-28-16.5c-10.1-4.3-15.1-8.6-15.1-12.9c-0.1-2.3,0.9-4.5,2.6-6.1  c1.7-1.6,3.9-2.4,6.2-2.4l0,0c5,0,9.3,3.5,12.7,10.5l21.6-11.5c-7-15-17.8-22.4-32.3-22.5l0,0c-6.9-0.1-13.7,1.8-19.6,5.4  c-2,1.2-3.8,2.7-5.5,4.2c-4.3,4.1-7.4,9.3-8.9,15.1c-0.8,3.1-1.2,6.2-1.2,9.4c0,7.4,2,13.5,5.9,18.3c3.1,3.7,8.4,7.4,16,11.2  c2.7,1.3,5.4,2.6,8,3.9c2.6,1.2,5.3,2.5,7.9,3.8c6.8,4,10.2,8.3,10.2,13c0,8.1-4.4,12.1-13.2,12.1l0,0c-4.2,0.1-8.3-1.5-11.4-4.4  c-1.6-1.5-4-4.9-7.1-10.4l-22.3,10.3c7.9,18.5,21.3,27.8,40.3,27.8l0,0c11.1,0,20.5-3.5,28.4-10.5c7.9-6.9,12.3-17,12.1-27.4  C272.1,73.3,269.7,66.5,264.9,61.5z"/>
                <path fill="#D15949" class="st0"
                      d="M128.5,30.8v42c0,17.1-6.7,25.6-20.2,25.6c-13.1,0-19.8-8-20.2-24.1V46.4c0-9.5-1-16.6-3-21.4  c-2.2-5.6-5.8-10.6-10.5-14.5C66.6,3.5,56.4,0,44,0C31.7,0,21.5,3.5,13.3,10.5C8.8,14.4,5.2,19.4,3.1,25C0.9,31.9-0.2,39.1,0,46.4  v73.3c0.6-0.7,1.3-1.4,2-2.1c6-5.7,13.7-9.4,21.9-10.6v-60c0-17.1,6.7-25.6,20.2-25.6c13.1,0,19.8,8,20.2,24.1v27.9  c0,9.5,1,16.6,3,21.4c2.2,5.6,5.8,10.6,10.4,14.5c8,7,18.2,10.5,30.6,10.5c12.4,0,22.6-3.5,30.7-10.5c4.6-3.9,8.1-8.9,10.3-14.5  c2.2-6.9,3.2-14.1,3.1-21.4V30.8H128.5z"/>
            </svg>
            </div>

            <div class="main-title text-xl">
              <small>/</small> <span class="font-semibold">Demo EHR</span>
            </div>
          </div>
        </div>

        <div class="px-3 mt-6">
          <div class="grid grid-cols-1">
            <router-link
                id="patients-menu-link"
                :to="{name: 'ehr.patients'}"
                class="menu-link"
                :class="{'menu-link-active': $route.name.indexOf('ehr.patient.') > -1}"
                active-class="menu-link-active">
              <div class="h-5 w-5 mr-3">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M0 0h24v24H0V0z" fill="none"/>
                  <path
                      d="M9 13.75c-2.34 0-7 1.17-7 3.5V19h14v-1.75c0-2.33-4.66-3.5-7-3.5zM4.34 17c.84-.58 2.87-1.25 4.66-1.25s3.82.67 4.66 1.25H4.34zM9 12c1.93 0 3.5-1.57 3.5-3.5S10.93 5 9 5 5.5 6.57 5.5 8.5 7.07 12 9 12zm0-5c.83 0 1.5.67 1.5 1.5S9.83 10 9 10s-1.5-.67-1.5-1.5S8.17 7 9 7zm7.04 6.81c1.16.84 1.96 1.96 1.96 3.44V19h4v-1.75c0-2.02-3.5-3.17-5.96-3.44zM15 12c1.93 0 3.5-1.57 3.5-3.5S16.93 5 15 5c-.54 0-1.04.13-1.5.35.63.89 1 1.98 1 3.15s-.37 2.26-1 3.15c.46.22.96.35 1.5.35z"/>
                </svg>
              </div>

              Patients
            </router-link>

            <router-link
                id="inbox-menu-link"
                :to="{name: 'ehr.inbox'}"
                class="menu-link"
                :class="{'menu-link-active': $route.name === 'ehr.transferRequest.show'}"
                active-class="menu-link-active">
              <inbox-badge :inbox="inboxInfo"/>
              Inbox
            </router-link>

            <router-link
                :to="{name: 'ehr.settings'}"
                class="menu-link"
                active-class="menu-link-active">
              <div class="w-5 h-5 mr-3">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M0 0h24v24H0V0z" fill="none"/>
                  <path
                      d="M19.43 12.98c.04-.32.07-.64.07-.98 0-.34-.03-.66-.07-.98l2.11-1.65c.19-.15.24-.42.12-.64l-2-3.46c-.09-.16-.26-.25-.44-.25-.06 0-.12.01-.17.03l-2.49 1c-.52-.4-1.08-.73-1.69-.98l-.38-2.65C14.46 2.18 14.25 2 14 2h-4c-.25 0-.46.18-.49.42l-.38 2.65c-.61.25-1.17.59-1.69.98l-2.49-1c-.06-.02-.12-.03-.18-.03-.17 0-.34.09-.43.25l-2 3.46c-.13.22-.07.49.12.64l2.11 1.65c-.04.32-.07.65-.07.98 0 .33.03.66.07.98l-2.11 1.65c-.19.15-.24.42-.12.64l2 3.46c.09.16.26.25.44.25.06 0 .12-.01.17-.03l2.49-1c.52.4 1.08.73 1.69.98l.38 2.65c.03.24.24.42.49.42h4c.25 0 .46-.18.49-.42l.38-2.65c.61-.25 1.17-.59 1.69-.98l2.49 1c.06.02.12.03.18.03.17 0 .34-.09.43-.25l2-3.46c.12-.22.07-.49-.12-.64l-2.11-1.65zm-1.98-1.71c.04.31.05.52.05.73 0 .21-.02.43-.05.73l-.14 1.13.89.7 1.08.84-.7 1.21-1.27-.51-1.04-.42-.9.68c-.43.32-.84.56-1.25.73l-1.06.43-.16 1.13-.2 1.35h-1.4l-.19-1.35-.16-1.13-1.06-.43c-.43-.18-.83-.41-1.23-.71l-.91-.7-1.06.43-1.27.51-.7-1.21 1.08-.84.89-.7-.14-1.13c-.03-.31-.05-.54-.05-.74s.02-.43.05-.73l.14-1.13-.89-.7-1.08-.84.7-1.21 1.27.51 1.04.42.9-.68c.43-.32.84-.56 1.25-.73l1.06-.43.16-1.13.2-1.35h1.39l.19 1.35.16 1.13 1.06.43c.43.18.83.41 1.23.71l.91.7 1.06-.43 1.27-.51.7 1.21-1.07.85-.89.7.14 1.13zM12 8c-2.21 0-4 1.79-4 4s1.79 4 4 4 4-1.79 4-4-1.79-4-4-4zm0 6c-1.1 0-2-.9-2-2s.9-2 2-2 2 .9 2 2-.9 2-2 2z"/>
                </svg>
              </div>

              Settings
            </router-link>

            <router-link
                :to="{name: 'ehr.signtest'}"
                class="menu-link"
                active-class="menu-link-active">
              <div class="w-5 h-5 mr-3">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M0 0h24v24H0V0z" fill="none"/>
                  <path
                      d="M19.43 12.98c.04-.32.07-.64.07-.98 0-.34-.03-.66-.07-.98l2.11-1.65c.19-.15.24-.42.12-.64l-2-3.46c-.09-.16-.26-.25-.44-.25-.06 0-.12.01-.17.03l-2.49 1c-.52-.4-1.08-.73-1.69-.98l-.38-2.65C14.46 2.18 14.25 2 14 2h-4c-.25 0-.46.18-.49.42l-.38 2.65c-.61.25-1.17.59-1.69.98l-2.49-1c-.06-.02-.12-.03-.18-.03-.17 0-.34.09-.43.25l-2 3.46c-.13.22-.07.49.12.64l2.11 1.65c-.04.32-.07.65-.07.98 0 .33.03.66.07.98l-2.11 1.65c-.19.15-.24.42-.12.64l2 3.46c.09.16.26.25.44.25.06 0 .12-.01.17-.03l2.49-1c.52.4 1.08.73 1.69.98l.38 2.65c.03.24.24.42.49.42h4c.25 0 .46-.18.49-.42l.38-2.65c.61-.25 1.17-.59 1.69-.98l2.49 1c.06.02.12.03.18.03.17 0 .34-.09.43-.25l2-3.46c.12-.22.07-.49-.12-.64l-2.11-1.65zm-1.98-1.71c.04.31.05.52.05.73 0 .21-.02.43-.05.73l-.14 1.13.89.7 1.08.84-.7 1.21-1.27-.51-1.04-.42-.9.68c-.43.32-.84.56-1.25.73l-1.06.43-.16 1.13-.2 1.35h-1.4l-.19-1.35-.16-1.13-1.06-.43c-.43-.18-.83-.41-1.23-.71l-.91-.7-1.06.43-1.27.51-.7-1.21 1.08-.84.89-.7-.14-1.13c-.03-.31-.05-.54-.05-.74s.02-.43.05-.73l.14-1.13-.89-.7-1.08-.84.7-1.21 1.27.51 1.04.42.9-.68c.43-.32.84-.56 1.25-.73l1.06-.43.16-1.13.2-1.35h1.39l.19 1.35.16 1.13 1.06.43c.43.18.83.41 1.23.71l.91.7 1.06-.43 1.27-.51.7 1.21-1.07.85-.89.7.14 1.13zM12 8c-2.21 0-4 1.79-4 4s1.79 4 4 4 4-1.79 4-4-1.79-4-4-4zm0 6c-1.1 0-2-.9-2-2s.9-2 2-2 2 .9 2 2-.9 2-2 2z"/>
                </svg>
              </div>

              Sign test
            </router-link>
          </div>
        </div>

      </div>

      <div class="px-3 py-4">
        <div class="grid grid-cols-1">
          <router-link
              id="logout-menu-link"
              to="/logout"
              active-class="menu-link-active"
              class="menu-link">
            <div class="w-5 h-5 mr-3">
              <svg viewBox="0 0 24 24" fill="currentColor">
                <g>
                  <path d="M0,0h24v24H0V0z" fill="none"/>
                </g>
                <g>
                  <path
                      d="M17,8l-1.41,1.41L17.17,11H9v2h8.17l-1.58,1.58L17,16l4-4L17,8z M5,5h7V3H5C3.9,3,3,3.9,3,5v14c0,1.1,0.9,2,2,2h7v-2H5V5z"/>
                </g>
              </svg>
            </div>

            Logout
          </router-link>
        </div>
      </div>
    </nav>

    <main class="ml-72 mb-14 w-full">
      <div v-if="organizationName" class="absolute top-6 right-12">
        <div class="flex flex-row items-center">
          <div
              class="rounded-full mb-1 mx-3 bg-white shadow-sm h-10 w-10 flex items-center justify-center overflow-hidden">
            <svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 0 24 24" width="24px" fill="#000000">
              <path d="M0 0h24v24H0V0z" fill="none"/>
              <path
                  d="M12 6c1.1 0 2 .9 2 2s-.9 2-2 2-2-.9-2-2 .9-2 2-2m0 10c2.7 0 5.8 1.29 6 2H6c.23-.72 3.31-2 6-2m0-12C9.79 4 8 5.79 8 8s1.79 4 4 4 4-1.79 4-4-1.79-4-4-4zm0 10c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z"/>
            </svg>
          </div>

          <div class="grid grid-cols-1 text-gray-800 text-sm font-semibold" id="current-organization-name">
            {{userIdentifier}} ({{ organizationName }})
          </div>
        </div>
      </div>

      <div class="w-full m-auto max-w-screen-2xl">
        <!-- Main content -->
        <div>
          <status-reporter type="error" :message="errorMessage.value"></status-reporter>
          <status-reporter type="info" :message="statusMessage.value"></status-reporter>
          <router-view/>
        </div>

        <status-bar :statusMessage="statusMessage"></status-bar>
      </div>
    </main>
  </div>
</template>

<script>
import {mapState} from "vuex";

import StatusBar from "../components/StatusBar.vue"
import StatusReporter from "../components/StatusReporter.vue"
import InboxBadge from "../components/InboxBadge.vue"

export default {
  components: {StatusBar, StatusReporter, InboxBadge},
  data() {
    return {
      inboxInfo: null,
    }
  },
  methods: {
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
    },
    tokenPayload() {
      const session = localStorage.getItem("session") || null;
      if (!session) {
        return null;
      }
      const [_header, rawPayload, _signature] = session.split('.');
      return JSON.parse(atob(rawPayload));
    },
    organizationName() {
      let p = this.tokenPayload;
      return p ? p.sub : null;
    },
    userIdentifier() {
      let p = this.tokenPayload;
      return p ? p.usi : null;
    },
    ...mapState(["statusMessage"])
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
