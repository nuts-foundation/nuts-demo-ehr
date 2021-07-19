<template>
  <div tabindex="-1"
       @focusin="open = true"
       @focusout="open = false"
  >
    <div class="border border-gray-300 rounded-md inline-flex items-center w-full px-2" >
      <input type="text"
             class="border-none h-8 p-0"
             :value="query"
             @input="query = $event.target.value; $emit('update:search', $event.target.value)"
      >
      <svg class="h-5 w-5 text-gray-400"
           @click="open = !open"
           xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 20 20" stroke="currentColor">
        <path fill-rule="evenodd"
              d="M10 3a1 1 0 01.707.293l3 3a1 1 0 01-1.414 1.414L10 5.414 7.707 7.707a1 1 0 01-1.414-1.414l3-3A1 1 0 0110 3zm-3.707 9.293a1 1 0 011.414 0L10 14.586l2.293-2.293a1 1 0 011.414 1.414l-3 3a1 1 0 01-1.414 0l-3-3a1 1 0 010-1.414z"
              clip-rule="evenodd"></path>
      </svg>

    </div>
    <ul v-if="open"
        class="border border-gray-300 rounded-md px-2 mt-1">
      <li class="hover:bg-gray-200 cursor-pointer"
          v-for="item in items"
          @click="select(item)"
      >
        <slot :item="item"></slot>
      </li>
      <li v-if="!items.length">
        No results
      </li>
    </ul>
  </div>
</template>

<script>
export default {
  name: 'auto-complete',
  data() {
    return {
      open: false,
      query: "",
    }
  },
  methods: {
    select(item) {
      this.$emit('update:selected', item)
      this.open = false
    },
    emitSearch() {
      this.$emit('update:search', this.query)
      return true
    }
  },
  props: {
    selected: Object,
    items: {
      type: Array,
      default: []
    },
    itemText: String
  },
  emits: ['update:search', 'update:selected']
}
</script>