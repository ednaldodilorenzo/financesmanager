export default {
  namespaced: true,
  state() {
    return {
      user: null,
    };
  },
  actions: {
    setUser({ commit }, payload) {
      commit("setUser", payload);
    },
  },
  mutations: {
    setUser(state, payload) {
      state.user = payload;
    },
  },
  getters: {
    getUser: (state) => state.user,
    getUserToken: (state) => state.user?.token,
    isAuthenticated: (state) => (state.user?.token ? true : false),
  },
};
