<template>
  <div class="px-12 py-8">
    <h1 class="mb-6 mt-12">Inbox</h1>

    <div class="bg-white p-5 shadow-lg rounded-lg">
      <table class="min-w-full divide-y divide-gray-200">
        <thead>
        <tr>
          <th></th>
          <th>Subject</th>
          <th>Status</th>
          <th>Sender</th>
          <th>Date</th>
        </tr>
        </thead>
        <tbody>
        <router-link
            :to="{name: 'ehr.transferRequest.show', params: {requestorDID: item.sender.did, fhirTaskID: item.resourceID}}"
            style="display: contents;"
            v-for="item in items"
        >
          <tr v-bind:class="{ 'hover:bg-gray-50': true, 'cursor-pointer': true, 'none': item.requiresAttention }">
            <td>
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2 inline" fill="none" viewBox="0 0 24 24"
                   stroke="currentColor" v-if="item.requiresAttention">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/>
              </svg>
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2 inline" fill="none" viewBox="0 0 24 24"
                   stroke="currentColor" v-if="!item.requiresAttention">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M3 19v-8.93a2 2 0 01.89-1.664l7-4.666a2 2 0 012.22 0l7 4.666A2 2 0 0121 10.07V19M3 19a2 2 0 002 2h14a2 2 0 002-2M3 19l6.75-4.5M21 19l-6.75-4.5M3 10l6.75 4.5M21 10l-6.75 4.5m0 0l-1.14.76a2 2 0 01-2.22 0l-1.14-.76"/>
              </svg>
            </td>
            <td>{{ item.title }}</td>
            <td>
              <transfer-status :status="item.status"/>
            </td>
            <td>{{ item.sender.name }}, {{
                item.sender.city
              }}
            </td>
            <td>{{ item.date }}</td>
          </tr>
        </router-link>
        </tbody>
      </table>
    </div>
  </div>
</template>
<script>
import TransferStatus from "../../components/TransferStatus.vue";

export default {
  components: {TransferStatus},
  data() {
    return {
      items: [],
    }
  },
  created() {
    this.fetchData()
  },
  methods: {
    fetchData() {
      this.$api.getInbox()
          .then((response) => this.items = response)
          .catch(error => this.$status.error(error))
    }
  }
}
</script>
