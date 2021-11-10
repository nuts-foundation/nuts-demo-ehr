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

      <button
          @click="$router.push({name: 'ehr.patients.new'})"
          class="float-right inline-flex items-center bg-nuts w-10 h-10 rounded-lg justify-center shadow-md"
      >
        <svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 0 24 24" width="24px" fill="#fff">
          <path d="M0 0h24v24H0V0z" fill="none"/>
          <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/>
        </svg>
      </button>
    </div>
  </div>

  <div class="px-12">
    <div class="grid gap-5 xl:grid-cols-4 lg:grid-cols-3 md:grid-cols-2 sm:grid-cols-1">
      <div class="bg-white p-6 shadow-md rounded cursor-pointer hover:shadow-lg opacity-0 transition-opacity"
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
          .then((response) => this.patients = response)
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
