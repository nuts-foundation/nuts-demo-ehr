<template>
  <h1 class="page-title">Inbox</h1>

  <ol>
    <li v-for="item in items">{{ item.title }}</li>
  </ol>
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
