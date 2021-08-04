<!--
Loosely based on on the w3c aria example.
https://www.w3.org/TR/wai-aria-practices-1.2/examples/combobox/combobox-autocomplete-list.html
-->
<template>
  <div class="combobox combobox-list"
       role="combobox"
  >
    <div class="group">
      <input type="text"
             class="cb_edit"
             @input="updateSearch($event.target.value)"
             :aria-expanded="expanded"
             aria-autocomplete="list"
             @focusin="focus"
             @focusout="focusLost"
             @keyup.down="highlightNext"
             @keyup.up="highlightPrev"
             @keyup.enter="selectHighlighted"
      >
      <svg class="button"
           @click="expanded = !expanded"
           xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 20 20" stroke="currentColor">
        <path fill-rule="evenodd"
              d="M10 3a1 1 0 01.707.293l3 3a1 1 0 01-1.414 1.414L10 5.414 7.707 7.707a1 1 0 01-1.414-1.414l3-3A1 1 0 0110 3zm-3.707 9.293a1 1 0 011.414 0L10 14.586l2.293-2.293a1 1 0 011.414 1.414l-3 3a1 1 0 01-1.414 0l-3-3a1 1 0 010-1.414z"
              clip-rule="evenodd"></path>
      </svg>
    </div>
    <ul v-if="expanded"
        tabindex="-1"
        role="listbox">
      <li v-for="(item, idx) in items"
          @click="select(item)"
          role="option"
      >
        <div :class="{'bg-gray-100': highlighted === idx}">
          <slot :item="item" ></slot>
        </div>
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
      expanded: false,
      selected: null,
      highlighted: null,
    }
  },
  methods: {
    select(item) {
      this.$emit('selected', item)
      this.expanded = false
    },
    updateSearch(val) {
      this.highlighted = null
      this.expanded = val !== ""
      this.$emit('search', val)
    },
    focus() {
      this.expanded = this.items.length > 0
    },
    focusLost(event) {
      const target = event.relatedTarget
      // Find out if the event took place in the listbox or else on the page. If it was the listbox keep it expanded, otherwise contract.
      if (target === null || !("role" in target.attributes) || target.attributes.role.value !== "listbox") {
        this.expanded = false
      }
    },
    // The next 3 methods are for keyboard navigation: up, down, enter.
    highlightPrev() {
      if (this.items.length === 0) return
      if (this.highlighted == null || this.highlighted === 0) {
        this.highlighted = this.items.length - 1
      } else {
        this.highlighted = this.highlighted - 1 % this.items.length
      }
    },
    highlightNext() {
      if (this.items.length === 0) return
      if (this.highlighted == null) {
        this.highlighted = 0
      } else {
        this.highlighted = (this.highlighted + 1) % this.items.length
      }
      console.log("highlight next", this.highlighted)
    },
    selectHighlighted() {
      if (this.items.length === 0) return
      this.select(this.items[this.highlighted])
    }
  },
  props: {
    items: {
      type: Array,
      default: []
    },
    itemText: String
  },
  emits: ['search', 'selected']
}
</script>

<style>

.combobox:focus {
  @apply border border-red-200;
}

.combobox .group {
  @apply border border-gray-300 rounded-md inline-flex items-center w-full px-2;
}

.combobox .group input {
  @apply border-none h-8 p-0;
}

.combobox .group .button {
  @apply h-5 w-5 text-gray-400;
}

ul[role="listbox"] {
  @apply border border-gray-300 rounded-md px-2 mt-1;
}

li[role="option"] {
  @apply hover:bg-gray-200 cursor-pointer;
}

</style>