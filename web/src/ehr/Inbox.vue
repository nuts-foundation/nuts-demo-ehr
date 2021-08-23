<template>
  <div>
    <h1 class="page-title">Inbox</h1>

    <table class="min-w-full divide-y divide-gray-200 mt-2">
      <thead class="bg-gray-50">
      <tr>
        <th></th>
        <th>Subject</th>
        <th>Sender</th>
        <th>Date</th>
      </tr>
      </thead>
      <tbody>
      <tr class="hover:bg-gray-100 cursor-pointer"
          v-for="item in items"
          @click="$router.push({name: 'ehr.transferRequest.show', params: {requestorDID: item.sender.did, fhirTaskID: item.resourceID}})">
        <td>
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-2 inline" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/>
          </svg>
        </td>
        <td v-bind:class="{ 'font-extrabold': item.requiresAttention }">{{ item.title }}</td>
        <td v-bind:class="{ 'font-extrabold': item.requiresAttention }">{{ item.sender.name }}, {{ item.sender.city }}</td>
        <td v-bind:class="{ 'font-extrabold': item.requiresAttention }">{{ item.date }}</td>
      </tr>
      </tbody>
    </table>
  </div>
</template>
<script>
export default {
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
