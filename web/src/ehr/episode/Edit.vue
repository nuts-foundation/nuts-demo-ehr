<template>
  <div>
    <h1>View episode</h1>
    <episode-fields v-if="episode"
                    :episode="episode"
                    mode="edit"
    ></episode-fields>

    <div class="mt-6">
      <label>Status</label>
      <div v-if="episode"
           class="bg-white p-5 shadow-sm rounded-lg mb-3">
        {{ episode.status }}
      </div>
    </div>

    <div class="mt-6">
      <button
          class="float-right inline-flex items-center bg-nuts w-10 h-10 rounded-lg justify-center shadow-md"
          @click="$router.push({name: 'ehr.patient.episode.newReport'})"
      >
        <svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 0 24 24" width="24px" fill="#fff">
          <path d="M0 0h24v24H0V0z" fill="none"/>
          <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/>
        </svg>
      </button>
      <label>Reports</label>
      <div v-if="reports.length > 0"
           class="bg-white p-5 shadow-sm rounded-lg mb-3">
        <table class="reports-list min-w-full divide-y divide-gray-200">
          <thead>
          <tr>
            <th>Type</th>
            <th>Value</th>
            <th>Source</th>
          </tr>
          </thead>
          <tbody>
          <tr v-for="report in reports">
            <td>{{ report.type }}</td>
            <td>{{ truncate(report.value, 30) }}</td>
            <td>{{ report.source }}</td>
          </tr>
          </tbody>
        </table>
      </div>

    </div>

    <router-view></router-view>
  </div>
</template>
<script>
import EpisodeFields from "./EpisodeFields.vue";

export default {
  components: {EpisodeFields},
  data() {
    return {
      episode: null,
      reports: [],
    }
  },
  emits: ['statusUpdate'],
  methods: {
    truncate(str, n) {
      return (str.length > n) ? str.substr(0, n - 1) + '...' : str
    },
    fetchEpisode(episodeID) {
      this.$api.getEpisode({episodeID})
          .then(episode => this.episode = episode)
          .catch(e => this.$status.error(e))
    },

    fetchReports(patientID, episodeID) {
      this.$api.getReports({patientID, episodeID})
          .then(reports => this.reports = reports)
          .catch(e => this.$status.error(e))
    }

  },
  mounted() {
    this.fetchEpisode(this.$route.params.episodeID)
    this.fetchReports(this.$route.params.id, this.$route.params.episodeID)
  }
}
</script>