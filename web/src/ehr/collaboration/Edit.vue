<template>
  <div>
    <div v-if="collaboration">
      Status: {{ collaboration.status }}
    </div>
  </div>
</template>
<script>
export default {
  data() {
    return {
      collaboration: {}
    }
  },
  emits: ['statusUpdate'],
  methods: {
    fetchCollaboration(collaborationID) {
      this.$api.getCollaboration({collaborationID})
          .then(collaboration => this.collaboration = collaboration)
          .catch(e => this.$status.error(e))
    }
  },
  mounted() {
    this.fetchCollaboration(this.$route.params.collaborationID)
  }
}
</script>