<template>
  <modal-window
      title="Create new Dossier report"
      type="add"
      :confirm-fn="submit"
      confirm-text="Create Report"
      :cancel-route="{name: 'ehr.patient.episode.edit', params: {id: $route.params.id, episodeID: $route.params.episodeID}}">
    <div class="mt-4">
      <h1>New report</h1>
      <form>
        <form-errors-banner :errors="formErrors">
        </form-errors-banner>
        <label>Heart rate
          <input type="text" v-model="report.heartRate">
        </label>
      </form>
    </div>
  </modal-window>

</template>

<script>
import ModalWindow from "../../../components/ModalWindow.vue";
import FormErrorsBanner from "../../../components/FormErrorsBanner.vue"

export default {
  name: "NewReport",
  components: {ModalWindow, FormErrorsBanner},
  data() {
    return {
      formErrors: [],
      report: {
        heartRate: null,
      }
    }
  },
  methods: {
    checkForm(e) {
      // reset the errors
      this.formErrors.length = 0

      if (!this.report.heartRate || this.report.heartRate < 1) {
        this.formErrors.push("The heart rate should be a positive number")
      }

      return this.formErrors.length === 0
    },
    submit() {
      if (!this.checkForm()) {
        return false
      }
      this.loading = true;

      let patientID = this.$route.params.id
      const payload = {
        type: "heartRate",
        patientID,
        value: this.report.heartRate.toString(),
        episodeID: this.$route.params.episodeID,
      };

      this.$api.createReport({
        body: payload,
        patientID,
      })
      this.$router.push({name: 'ehr.patient.episode.edit', params: {id: this.$route.params.id}})
    },
  },
}
</script>
