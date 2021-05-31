<template>
  <div class="flex">
    <!-- Navigation -->
    <nav class="min-h-screen sticky top-0 w-96 border border-r-1 border-gray-200">
      <div class="flex justify-center pt-6">
      </div>
      <h1 class="mt-4 mb-12 text-2xl text-center">Nuts EHR</h1>
      <div class="grid grid-cols-1">
        <router-link
            :to="{name: 'ehr.patients'}"
            class="menu-link">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-2" fill="none" viewBox="0 0 24 24"
               stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"/>
          </svg>
          Patients
        </router-link>
        <router-link
            to="/logout"
            active-class="text-blue-400"
            class="menu-link">
          Logout
        </router-link>
      </div>

    </nav>

    <main class="relative w-full p-5">
      <!-- Main content -->
      <div class="max-w-4xl mx-auto my-8">
        <router-view @statusUpdate="updateStatus"></router-view>
      </div>
      <status-bar :statusMessage="eventMessage"></status-bar>
    </main>
  </div>
</template>

<script>
import StatusBar from "../components/StatusBar.vue";

export default {
  components: {StatusBar},
  data() {
    return {
      eventMessage: ''
    }
  },
  methods: {
    updateStatus(status) {
      this.eventMessage = status
    }
  }
}
</script>

<style>
.menu-link {
  @apply px-5 h-14 inline-flex items-center text-sm hover:bg-blue-50;
}
</style>
