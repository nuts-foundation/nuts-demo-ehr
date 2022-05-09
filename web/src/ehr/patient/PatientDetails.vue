<template>
  <div class="flex flex-row">
    <div class="w-24 h-24 flex-shrink-0 rounded-full bg-gray-300 overflow-hidden mr-5 opacity-0 transition-opacity"
         :class="{'opacity-100': Object.keys(patient).length > 0}">
      <avatar :gender="patient.gender" :avatar_url="patient.avatar_url"/>
    </div>

    <div class="w-full">
      <div class="flex justify-between">
        <router-link
            :to="{name: 'ehr.patient.overview', params: {id: patient.ObjectID}}"
            class="text-2xl mb-2 mr-4 hover:cursor-pointer hover:underline"
            v-if="editable && (patient.surname || patient.firstName)"
        >
          {{ patient.firstName }} {{ patient.surname }}
        </router-link>
        <p class="text-2xl mb-2 mr-4" v-else-if="!editable && (patient.surname || patient.firstName)"  id="patient-name-label">
          {{ patient.firstName }} {{ patient.surname }}
        </p>
        <div v-else-if="Object.keys(patient).length === 0">...</div>
        <div v-else class="text-2xl  mb-2 mr-4">Unknown patient</div>

        <button
            v-if="editable && $route.name !== 'ehr.patient.edit'"
            @click="$router.push({name: 'ehr.patient.edit', params: {id: patient.ObjectID}})"
            class="float-right inline-flex items-center bg-nuts w-10 h-10 rounded-lg justify-center shadow-md"
        >
          <svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 0 24 24" width="24px" fill="#fff">
            <path d="M0 0h24v24H0V0z" fill="none"/>
            <path
                d="M14.06 9.02l.92.92L5.92 19H5v-.92l9.06-9.06M17.66 3c-.25 0-.51.1-.7.29l-1.83 1.83 3.75 3.75 1.83-1.83c.39-.39.39-1.02 0-1.41l-2.34-2.34c-.2-.2-.45-.29-.71-.29zm-3.6 3.19L3 17.25V21h3.75L17.81 9.94l-3.75-3.75z"/>
          </svg>
        </button>
      </div>

      <div class="grid grid-cols-5 gap-x-6">
        <div v-if="patient.ssn">
          <div class="text-sm font-semibold">SSN</div>
          <span id="patient-ssn-label">{{ patient.ssn }}</span>
        </div>

        <div>
          <div class="text-sm font-semibold">Gender</div>
          <span id="patient-gender-label">{{ patient.gender ? patient.gender : (Object.keys(patient).length === 0 ? '...' : 'Unknown') }}</span>
        </div>

        <div>
          <div class="text-sm font-semibold">Birth date</div>
          <span id="patient-dob-label">{{ patient.dob ? patient.dob : (Object.keys(patient).length === 0 ? '...' : 'Unknown') }}</span>
        </div>

        <div v-if="patient.email">
          <div class="text-sm font-semibold">E-mail</div>
          <span id="patient-email-label">{{ patient.email }}</span>
        </div>

        <div v-if="patient.zipcode">
          <div class="text-sm font-semibold">Zipcode</div>
          <span id="patient-zipcode-label">{{ patient.zipcode }}</span>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
import Avatar from "../../components/Avatar.vue";

export default {
  components: {
    Avatar
  },
  props: {
    patient: Object,
    editable: {
      type: Boolean,
      default: false
    }
  },
}
</script>
