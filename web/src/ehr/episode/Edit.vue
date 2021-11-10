<template>
  <div>
    <h1>View episode</h1>
    <episode-fields v-if="episode"
                    :episode="episode"
                    mode="edit"
    ></episode-fields>

    <div v-if="episode"
         class="bg-white p-5 shadow-sm rounded-lg mb-3">
      Status: {{ episode.status }}
    </div>
  </div>
</template>
<script>
import EpisodeFields from "./EpisodeFields.vue";

export default {
  components: {EpisodeFields},
  data() {
    return {
      episode: null
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