import Vue from 'vue';

export const state = () => ({
    token: null,
    user: null
});

export const actions = {
    me({ commit }) {
        return this.$axios.get('/api/me').then(res => {
            commit('setUser', res.data);
        })
    }
};

export const mutations = {
    setToken(state, token) {
        state.token = token;
    },
    setUser(state, user) {
        state.user = user;
    }
}
