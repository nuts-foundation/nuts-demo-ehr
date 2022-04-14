import {createStore} from "vuex";

export default createStore({
  state() {
    return {
      statusMessage: ""
    }
  },
  mutations: {
    statusUpdate(state, message) {
      state.statusMessage = message
    }
  }
})