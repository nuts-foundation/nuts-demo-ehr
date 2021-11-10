<template>
  <div>
    <div v-if="episode">
      Status: {{ episode.status }}
    </div>
  </div>
</template>
<script>
export default {
  data() {
    return {
      episode: {}
    }
  },
  emits: ['statusUpdate'],
  methods: {
    fetchEpisode(episodeID) {
      this.$api.getEpisode({episodeID})
          .then(episode => this.episode = episode)
          .catch(e => this.$status.error(e))
    }
  },
  mounted() {
    this.fetchEpisode(this.$route.params.episodeID)
  }
}
</script>