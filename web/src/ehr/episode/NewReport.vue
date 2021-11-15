<template>
  <modal-window
      title="Create new report"
      type="add"
      :confirm-fn="submit"
      confirm-text="Create Report"
      :cancel-fn="{cancel}">
    <div class="mt-4">
      <form>
        <form-errors-banner :errors="formErrors"/>

        <label>Heart rate</label>
        <input type="text" v-model="report.heartRate">
      </form>
    </div>
  </modal-window>

</template>

<script>
import ModalWindow from "../../components/ModalWindow.vue";
import FormErrorsBanner from "../../components/FormErrorsBanner.vue"

export default {
  name: "NewReport",
  components: {ModalWindow, FormErrorsBanner},
  emits: ["added", "cancelled"],
  props: {
    episodeId: {
      type: String,
      required: true
    },
    patientId: {
      type: String,
      required: true
    }
  },
  data() {
    return {
      formErrors: [],
      report: {
        heartRate: null,
      }
    }
  },
  methods: {
    checkForm() {
      // reset the errors
      this.formErrors.length = 0

      if (!this.report.heartRate || this.report.heartRate < 1) {
        this.formErrors.push("The heart rate should be a positive number")
      }

      return this.formErrors.length === 0
    },
    cancel() {
      this.$emit("cancelled")
    },
    submit() {
      if (!this.checkForm()) {
        return false
      }

      this.loading = true

      const body = {
        type: "heartRate",
        patientID: this.patientId,
        value: this.report.heartRate.toString(),
        episodeID: this.episodeId,
      };

      this.$api.createReport({
        body,
        patientID: this.patientId,
      }).then(() => {
        this.$emit("added")
      })
    },
  },
}
</script>
