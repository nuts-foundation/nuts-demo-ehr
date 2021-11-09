<template>
  <modal-window
      title="Start a new dossier"
      type="add"
      :confirm-text="null"
      :cancel-route="{name: 'ehr.patient', params: {id: this.$route.params.id}}">
    <div class="mt-2">
      Do you want to transfer this patient to another care organization?
      You can make a set of predetermined ZIBs available to 1 specific care organization.

      <button
          @click.prevent="() => $router.push({name: 'ehr.patient.transfer.new', params: {id: this.$route.params.id}})"
          class="btn btn-primary block mt-2">
        Create transfer
      </button>
    </div>

    <div class="mt-6">
      Do you want to collaborate with a group of care professionals on a specific case?
      With the dossier type you invite other care professionals to exchange reports on a specific case.
      <button
          @click.prevent="() => $router.push({name: 'ehr.patient.collaboration.new', params: {id: this.$route.params.id}})"
          class="btn btn-primary block mt-2">
        Create collaboration
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
