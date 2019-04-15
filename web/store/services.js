import Service from '@/models/Service';

export const state = () => ({
    list: [],
    total: 0
});

export const actions = {
    fetch(ctx, { features = null, pagination }) {
        const q = `?f=${features}&p=${pagination.page}&l=${pagination.limit}&q=${pagination.search}`;
        return this.$axios.get('/api/services' + q)
            .then((res) => {
                ctx.commit('SET', res.data.services.map(w => new Service(w)));
                ctx.commit('TOTAL', res.data.total);
            });
    },
    delete(ctx, id) {
        return this.$axios.get(`/api/service/${id}/delete`)
            .then(() => ctx.commit('REMOVE', id));
    }
};

export const mutations = {
    SET(state, services) {
        state.list = services;
    },
    REMOVE(state, id) {
        const target = state.list.find(i => i.id === id);
        if (target) {
            state.list.splice(state.list.indexOf(target), 1);
        }
    },
    TOTAL(state, total) {
        state.total = total;
    }
};
