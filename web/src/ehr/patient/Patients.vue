<template>
  <div class="px-12 py-8">
    <form @submit.prevent="list" class="inline-flex mb-10">
      <button class="btn btn-secondary px-1 bg-transparent shadow-none">
        <svg width="18" height="18" viewBox="0 0 18 18" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path
              d="M12.5 11H11.71L11.43 10.73C12.63 9.33002 13.25 7.42002 12.91 5.39002C12.44 2.61002 10.12 0.390015 7.32 0.0500152C3.09 -0.469985 -0.47 3.09001 0.05 7.32001C
0.39 10.12 2.61 12.44 5.39 12.91C7.42 13.25 9.33 12.63 10.73 11.43L11 11.71V12.5L15.25 16.75C15.66 17.16 16.33 17.16 16.74 16.75C17.15 16.34 17.15 15.67 16.74 15.26L12.5 11ZM6.5
 11C4.01 11 2 8.99002 2 6.50002C2 4.01002 4.01 2.00002 6.5 2.00002C8.99 2.00002 11 4.01002 11 6.50002C11 8.99002 8.99 11 6.5 11Z"
              fill="#111"/>
        </svg>
      </button>

      <input class="bg-transparent border-0 shadow-none w-full hover:border-0" placeholder="Search.." type="text" v-model="query">
    </form>

    <div class="flex justify-between items-center mb-2">
      <h1>Patients</h1>

      <div class="float-right inline-flex">
        <button
            id="remote-patient-button"
            @click="$router.push({name: 'ehr.patients.remote'})"
            class="items-center  inline-flex bg-nuts w-10 h-10 rounded-lg justify-center shadow-md mr-2">
          <svg xmlns="http://www.w3.org/2000/svg" fill="#fff" viewBox="0 0 24 24" width="24px" height="24px">
            <path stroke-linecap="round" stroke-linejoin="round"
                  d="M12.75 3.03v.568c0 .334.148.65.405.864l1.068.89c.442.369.535 1.01.216 1.49l-.51.766a2.25 2.25 0 0 1-1.161.886l-.143.048a1.107 1.107 0 0 0-.57 1.664c.369.555.169 1.307-.427 1.605L9 13.125l.423 1.059a.956.956 0 0 1-1.652.928l-.679-.906a1.125 1.125 0 0 0-1.906.172L4.5 15.75l-.612.153M12.75 3.031a9 9 0 0 0-8.862 12.872M12.75 3.031a9 9 0 0 1 6.69 14.036m0 0-.177-.529A2.25 2.25 0 0 0 17.128 15H16.5l-.324-.324a1.453 1.453 0 0 0-2.328.377l-.036.073a1.586 1.586 0 0 1-.982.816l-.99.282c-.55.157-.894.702-.8 1.267l.073.438c.08.474.49.821.97.821.846 0 1.598.542 1.865 1.345l.215.643m5.276-3.67a9.012 9.012 0 0 1-5.276 3.67m0 0a9 9 0 0 1-10.275-4.835M15.75 9c0 .896-.393 1.7-1.016 2.25"/>
          </svg>
        </button>
        <button
            id="new-patient-button"
            @click="$router.push({name: 'ehr.patients.new'})"
            class="items-center inline-flex bg-nuts w-10 h-10 rounded-lg justify-center shadow-md mr-2">
          <svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 0 24 24" width="24px" fill="#fff">
            <path d="M0 0h24v24H0V0z" fill="none"/>
            <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/>
          </svg>
        </button>
      </div>
    </div>
  </div>

  <div class="px-12">
    <div class="grid gap-5 xl:grid-cols-4 lg:grid-cols-3 md:grid-cols-2 sm:grid-cols-1">
      <div :data-patient-ssn="patient.ssn"
           class="bg-white p-6 shadow-md rounded cursor-pointer hover:shadow-lg opacity-0 transition-opacity"
           :class="{'opacity-100': state === 'done'}"
           v-for="(patient, i) in patients"
           @click="$router.push({name: 'ehr.patient', params: {id: patient.ObjectID}})"
           :style="{'transition-duration': `${.05 * (i+1)}s !important`}">
        <div class="inline-flex items-center">
          <div class="flex-shrink-0 w-11 h-11 mr-3 rounded-full bg-gray-300 overflow-hidden">
            <avatar
                :gender="patient.gender"
                :avatar_url="patient.avatar_url"
            />
          </div>

          <div>
            <h3 class="font-semibold text-gray-700 text-md">
              {{ patient.firstName }} {{ patient.surname }}
            </h3>

            <p class="text-gray-500 text-md" :title="patient.dob">{{ calculateAge(patient.dob) }} yr <small>/
              {{ patient.dob }}</small></p>
          </div>
        </div>
      </div>
    </div>

    <div v-if="state !== 'loading' && patients.length === 0 && query">No results</div>
  </div>

</template>
<script>
import Avatar from "../../components/Avatar.vue";

export default {
  components: {
    Avatar
  },
  data() {
    return {
      patients: [],
      query: '',
      state: 'initial',
    }
  },
  mounted() {
    this.list()
  },
  methods: {
    list() {
      this.patients = []
      this.state = 'loading';
      let params = {};
      if (this.query !== "") {
        params.name = this.query
      }
      this.$api.getPatients(params)
          .then((result) => this.patients = result.data)
          .catch(error => this.$status.error(error))
          .finally(() => {
            this.$nextTick(() => setTimeout(() => this.state = 'done', 10));
          })
    },

    calculateAge(dob) {
      const [y, m, d] = dob.split("-");
      const date = new Date(y, m, d);
      const now = new Date();

      return new Date(now - date).getUTCFullYear() - 1970;
    }
  },
}
</script>
