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

    <div class="mt-8">
      <div class="flex justify-between items-center mb-4">
        <h2>Reports</h2>

        <button
            class="float-right inline-flex items-center bg-nuts w-10 h-10 rounded-lg justify-center shadow-md"
            @click="$router.push({name: 'ehr.patient.episode.newReport'})"
        >
          <svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 0 24 24" width="24px" fill="#fff">
            <path d="M0 0h24v24H0V0z" fill="none"/>
            <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/>
          </svg>
        </button>
      </div>

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

    <div class="mt-8">
      <h2>Collaborations</h2>

      <div class="bg-white p-5 shadow-sm rounded-lg mt-4">
        <table class="min-w-full divide-y divide-gray-200">
          <thead>
          <tr>
            <th>Organization</th>
          </tr>
          </thead>
          <tbody>
            <tr v-for="collaboration in collaborations">
              <td>{{ collaboration.organizationName }}</td>
            </tr>

            <tr>
              <td>
                <auto-complete
                    :items="organizations"
                    @selected="chooseCollaboration"
                    @search="searchOrganizations"
                    v-slot="slotProps">
                  {{ slotProps.item.name }}
                </auto-complete>
              </td>
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
import AutoComplete from "../../components/Autocomplete.vue"

export default {
  components: {EpisodeFields, AutoComplete},
  data() {
    return {
      episode: null,
      reports: [],
      collaborations: [],
      organizations: [],
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
    },
    fetchCollaborations(episodeID) {
      this.$api.getCollaboration({episodeID})
          .then(collaborations => this.collaborations = collaborations)
          .catch(e => this.$status.error(e))
    },
    chooseCollaboration(collaboration) {
      const episodeID = this.$route.params.episodeID

      this.$api.createCollaboration({episodeID, body: {sender: {did: collaboration.did}}})
          .then(() => this.fetchCollaborations(episodeID))
          .catch(error => this.$status.error(error))
    },
    searchOrganizations(query) {
      this.$api.searchOrganizations({query: query, didServiceType: "zorginzage-demo"})
          .then((organizations) => {
            this.organizations = organizations
          })
          .catch(error => this.$status.error(error))
    },
  },
  mounted() {
    this.fetchEpisode(this.$route.params.episodeID)
    this.fetchReports(this.$route.params.id, this.$route.params.episodeID)
    this.fetchCollaborations(this.$route.params.episodeID)
  }
}
</script>