<template>
  <modal-window title="Start a new dossier" type="add" :cancel-route="{name: 'ehr.patient', params: {id: this.$route.params.id}}">
    <div class="mt-4">
      <span>
        Do you want to transfer this patient to another care organization?
        You can make a set of predetermined ZIBs available to 1 specific care organization.
      </span>
      <button class="btn btn-primary block mt-1"
              @click="$router.push({name: 'ehr.patient.transfer.new', params: {id: this.$route.params.id}})">
        Transfer
      </button>
    </div>

    <!--
    Not supported(yet)

    <div class="mt-4">
      Do you want to collaborate with a group of care professionals on a specific case?
      With the dossier type you invite other care professionals to exchange reports on a specific case.
      <button class="btn btn-primary block mt-1">
        Collaborate
      </button>
    </div>
    -->
  </modal-window>
</template>
<script>

import ModalWindow from "../../../components/ModalWindow.vue";

export default {
  components: {ModalWindow},
}
</script>
