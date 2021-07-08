<template>
  <div class="m-4">
    <h1>Autocomplete</h1>
    query: {{ search }}
    <auto-complete
        :items="this.items"
        v-model:searchQuery="search"
        v-model:selected="selectedCat"
        v-slot="slotProps"
    >
       {{slotProps.item.id}}
    </auto-complete>
    <img class="h-full" :src="selectedCat.url" alt="selected cat"/>
  </div>
</template>
<script>
// This file can be used to develop and test components with different configurations.

import AutoComplete from "./components/Autocomplete.vue";

export default {
  components: {AutoComplete},
  data: () => ({
    descriptionLimit: 60,
    apiResults: [],
    search: null,
    isLoading: false,
    selectedCat: {}
  }),
  watch: {
    search(val) {
      if (val === "") {
        this.apiResults = []
        return
      }
      // Items have already been requested
      if (this.isLoading) return

      this.isLoading = true

      // Lazily load input items
      fetch(`https://cataas.com/api/cats?tags=${val}&limit=6`)
          .then(res => res.json())
          .then(res => {
            this.apiResults = res
          })
          .catch(err => {
            console.log(err)
          })
          .finally(() => (this.isLoading = false))
    },
  },
  computed: {
    items() {
      return this.apiResults.map(entry => {
        return {id: entry.id, url: `https://cataas.com/cat/${entry.id}`}
      });
    },
  }

}
</script>