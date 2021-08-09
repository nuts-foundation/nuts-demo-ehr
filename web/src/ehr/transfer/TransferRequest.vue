<template>
  <div>
    <h1 class="page-title">Transfer Request</h1>

    <div class="mt-4" v-if="transfer">
      <div class="bg-gray-50 font-bold">Requesting Care Organization</div>
      <div>
        {{ transfer.sender.name }}, {{ transfer.sender.city }}
      </div>
    </div>

    <div class="mt-4" v-if="transfer">
      <div class="bg-gray-50 font-bold">Transfer date</div>
      <div>
        {{ transfer.transferDate }}
      </div>
    </div>

    <div class="mt-4" v-if="transfer">
      <div class="bg-gray-50 font-bold">Condition</div>
      <div>
        {{ transfer.condition }}
      </div>
    </div>


  </div>
</template>
<script>
export default {
  data() {
    return {
      items: [],
    }
  },
  created() {
    this.fetchData()
  },
  methods: {
    fetchData() {
      this.$api.getInbox()
          .then((response) => this.items = response)
          .catch(error => this.$status.error(error))
    }
  }
}
</script>
