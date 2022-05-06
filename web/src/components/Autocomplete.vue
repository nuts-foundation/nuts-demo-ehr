<!--
Loosely based on on the w3c aria example.
https://www.w3.org/TR/wai-aria-practices-1.2/examples/combobox/combobox-autocomplete-list.html
-->
<template>
  <div class="combobox combobox-list"
       role="combobox"
  >
    <div class="group custom-select">
      <input id="transfer-receiver-input"
             type="text"
             class="cb_edit px-4 py-2"
             @input="updateSearch($event.target.value)"
             :aria-expanded="expanded"
             aria-autocomplete="list"
             @focusin="focus"
             @focusout="focusLost"
             @keyup.down="highlightNext"
             @keyup.up="highlightPrev"
             @keyup.enter="selectHighlighted"
      >

      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="#444" @click="expanded = !expanded">
        <path d="M24 24H0V0h24v24z" fill="none" opacity=".87"/>
        <path d="M16.59 8.59L12 13.17 7.41 8.59 6 10l6 6 6-6-1.41-1.41z"/>
      </svg>
    </div>
    <ul v-if="expanded"
        tabindex="-1"
        role="listbox"
        class="rounded-b-lg border overflow-hidden">
      <li v-for="(item, idx) in items"
          @click="select(item)"
          role="option"
          class="bg-white p-3 cursor-pointer hover:bg-blue-100"
      >
        <div :class="{'bg-gray-100': highlighted === idx}">
          <slot :item="item"></slot>
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
