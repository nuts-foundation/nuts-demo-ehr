<template>
  <div class="flex flex-row">
    <div class="w-24 h-24 flex-shrink-0 rounded-full bg-gray-300 overflow-hidden mr-5">
      <avatar :gender="patient.gender" :avatar_url="patient.avatar_url"/>
    </div>

    <div class="w-full">
      <div class="flex justify-between">
        <router-link
            :to="{name: 'ehr.patient.overview', params: {id: patient.ObjectID}}"
            class="text-2xl mb-2 mr-4 hover:cursor-pointer hover:underline"
            v-if="patient.surname || patient.firstName"
        >{{ patient.surname }}, {{ patient.firstName }}
        </router-link>

        <h1 v-else class="text-2xl mb-2 mr-4">Unknown patient</h1>

        <button
            @click="$router.push({name: 'ehr.patient.edit', params: {id: patient.ObjectID}})"
            class="float-right inline-flex items-center bg-blue-700 w-10 h-10 rounded-lg justify-center shadow-md"
        >
          <svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 0 24 24" width="24px" fill="#fff">
            <path d="M0 0h24v24H0V0z" fill="none"/>
            <path
                d="M14.06 9.02l.92.92L5.92 19H5v-.92l9.06-9.06M17.66 3c-.25 0-.51.1-.7.29l-1.83 1.83 3.75 3.75 1.83-1.83c.39-.39.39-1.02 0-1.41l-2.34-2.34c-.2-.2-.45-.29-.71-.29zm-3.6 3.19L3 17.25V21h3.75L17.81 9.94l-3.75-3.75z"/>
          </svg>
        </button>
      </div>

      <div class="grid grid-cols-2 gap-x-6">
        <div v-if="patient.ssn">
          <div class="text-sm font-semibold">SSN</div>
          {{ patient.ssn }}
        </div>
        <div>
          <div class="text-sm font-semibold">Gender</div>
          {{ patient.gender }}
        </div>
        <div>
          <div class="text-sm font-semibold">Birth date</div>
          {{ patient.dob ? patient.dob : "unknown" }}
        </div>
        <div v-if="patient.email">
          <div class="text-sm font-semibold">E-mail</div>
          {{ patient.email }}
        </div>
        <div v-if="patient.zipcode">
          <div class="text-sm font-semibold">Zipcode</div>
          {{ patient.zipcode }}
        </div>
      </div>
    </div>
  </div>
</template>
<script>
import Avatar from "../../components/Avator.vue";

export default {
  components: {
    Avatar
  },
  props: {
    patient: Object,
  },
}
</script>
