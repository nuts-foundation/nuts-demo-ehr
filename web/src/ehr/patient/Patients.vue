<template>
  <div class="px-12 py-8">
    <h1 class="mb-4">Patients</h1>

    <form @submit.prevent="list" class="inline-flex">
      <button class="btn btn-secondary px-1 bg-transparent shadow-none">
        <svg width="18" height="18" viewBox="0 0 18 18" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path
              d="M12.5 11H11.71L11.43 10.73C12.63 9.33002 13.25 7.42002 12.91 5.39002C12.44 2.61002 10.12 0.390015 7.32 0.0500152C3.09 -0.469985 -0.47 3.09001 0.05 7.32001C
0.39 10.12 2.61 12.44 5.39 12.91C7.42 13.25 9.33 12.63 10.73 11.43L11 11.71V12.5L15.25 16.75C15.66 17.16 16.33 17.16 16.74 16.75C17.15 16.34 17.15 15.67 16.74 15.26L12.5 11ZM6.5
 11C4.01 11 2 8.99002 2 6.50002C2 4.01002 4.01 2.00002 6.5 2.00002C8.99 2.00002 11 4.01002 11 6.50002C11 8.99002 8.99 11 6.5 11Z"
              fill="#111"/>
        </svg>
      </button>

      <input class="bg-transparent border-0 w-48" placeholder="Search.." type="text" v-model="query">
    </form>

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

  <div class="px-12 py-8">
    <div class="grid gap-5 grid-cols-4">
      <div class="bg-white p-6 shadow-md rounded cursor-pointer hover:shadow-lg"
           v-for="patient in patients"
           @click="$router.push({name: 'ehr.patient', params: {id: patient.ObjectID}})">
        <div class="inline-flex items-center mb-6">
          <div class="flex-shrink-0 w-11 h-11 mr-3 rounded-full bg-gray-300 overflow-hidden">
            <avatar
                :gender="patient.gender"
                :avatar_url="patient.avatar_url"
            />
          </div>

          <h3 class="font-bold text-gray-900 text-md">
            {{ patient.firstName }} {{ patient.surname }}
          </h3>
        </div>

        <h5 class="font-semibold text-sm">Gender</h5>
        <div>{{ patient.gender }}</div>

        <h5 class="font-semibold text-sm mt-2">Age</h5>
        <p :title="patient.dob">{{ calculateAge(patient.dob) }} <small class="text-gray-500">/ {{ patient.dob }}</small>
        </p>
      </div>
    </div>

    <div class="text-gray-500" v-if="loading">Loading...</div>
    <div v-if="!loading && patients.length == 0 && query">No results</div>
  </div>

</template>
<script>
import Avatar from "../../components/Avator.vue";

export default {
  components: {
    Avatar
  },
  data() {
    return {
      patients: [],
      query: "",
      loading: false,
    }
  },
  mounted() {
    this.list()
  },
  methods: {
    list() {
      this.patients = []
      this.loading = true
      let params = {};
      if (this.query != "") {
        params.name = this.query
      }
      this.$api.getPatients(params)
          .then((response) => this.patients = response)
          .catch(error => this.$status.error(error))
          .finally(() => this.loading = false)
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
