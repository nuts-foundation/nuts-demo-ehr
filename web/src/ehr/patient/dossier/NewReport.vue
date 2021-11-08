<template>
  <modal-window
      title="Create new Dossier report"
      type="add"
      :confirm-fn="submit"
      confirm-text="Create Report"
      :cancel-route="{name: 'ehr.patient', params: {id: this.$route.params.id}}">
    <div class="mt-4">
      <h1>New report</h1>
      <form>
        <div class="sticky top-0 z-10 p-3 bg-red-100 text-red-500 rounded-md" v-if="formErrors.length">
          <label class="text-red-500">Please correct the following error{{
              formErrors.length === 0 ? '' : 's'
            }}:</label>

          <ul class="text-sm">
            <li v-for="error in formErrors">— {{ error }}</li>
          </ul>
        </div>
        <label>Heart rate
          <input type="text" v-model="report.heartRate">
        </label>
      </form>
    </div>
  </modal-window>

</template>

<script>
import ModalWindow from "../../../components/ModalWindow.vue";

export default {
  name: "NewReport",
  components: {ModalWindow},
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
        patientID: patientID,
        value: this.report.heartRate.toString(),
      };

      this.$api.createReport({
        body: payload,
        patientID: patientID,
      })
      this.$router.push({name: 'ehr.patient', params: {id: this.$route.params.id}})
    },
  },
}
</script>
