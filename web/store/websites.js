import Website from '@/models/Website';

export const state = () => ({
    list: []
});

export const getters = {

};

export const actions = {
    fetch(ctx, filters = null) {
        let q = '';
        if (filters) {
            q = '?f=' + filters;
        }
        return this.$axios.get('/api/websites' + q)
            .then(res => ctx.commit('SET', res.data.map(w => new Website(w))));
    }
};

export const mutations = {
    SET(state, websites) {
        state.list = websites;
    },
    ADD(state, websites) {
        state.list = state.list.concat(websites);
    },
    SET_LOADING(state, { id, status }) {
        const target = state.list.find(i => i.id === id);
        if (target) {
            target.loading = status;
        }
    },
    UPDATE(state, website) {
        const target = state.list.find(i => i.id === website.id);
        if (target) {
            Object.assign(target, website);
        }
    }
};
