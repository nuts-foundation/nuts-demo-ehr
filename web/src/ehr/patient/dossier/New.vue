<template>
  <modal-window
      title="Start a new dossier"
      type="add"
      :confirm-text="null"
      :cancel-route="{name: 'ehr.patient', params: {id: this.$route.params.id}}">
    <div class="mt-2">
      Do you want to create a new shared care plan for this patient, in which other care organization can participate?

      <button
          id="new-careplan-button"
          @click.prevent="() => $router.push({name: 'ehr.patient.careplan.new', params: {id: this.$route.params.id}})"
          class="btn btn-primary block mt-2">
        New shared care plan
      </button>
    </div>

    <div class="mt-2">
      Do you want to transfer this patient to another care organization?
      You can make a set of predetermined ZIBs available to 1 specific care organization.

      <button
          id="new-transfer-dossier-button"
          @click.prevent="() => $router.push({name: 'ehr.patient.transfer.new', params: {id: this.$route.params.id}})"
          class="btn btn-primary block mt-2">
        New transfer
      </button>
    </div>

    <div class="mt-6">
      Do you want to collaborate with a group of care professionals on a specific case?
      With the dossier type you invite other care professionals to exchange reports on a specific case.
      <button
          @click.prevent="() => $router.push({name: 'ehr.patient.episode.new', params: {id: this.$route.params.id}})"
          class="btn btn-primary block mt-2">
        New episode
      </button>
    </div>
  </modal-window>
</template>
<script>

import ModalWindow from "../../../components/ModalWindow.vue";

export default {
  components: {ModalWindow},
}
</script>
