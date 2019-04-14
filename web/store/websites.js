import Website from '@/models/Website';

export const state = () => ({
    list: [],
    total: 0,
    queued: 0,
    timeouted: 0
});

export const actions = {
    fetch(ctx, { filters = null, pagination }) {
        const q = `?f=${filters}&p=${pagination.page}&l=${pagination.limit}&q=${pagination.search}`;
        return this.$axios.get('/api/websites' + q)
            .then((res) => {
                ctx.commit('SET', res.data.websites.map(w => new Website(w)));
                ctx.commit('TOTAL', res.data.total);
            });
    },
    delete(ctx, ids) {
        return this.$axios.post('/api/websites/delete', ids)
            .then(res => ctx.commit('REMOVE', ids));
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
    },
    REMOVE(state, ids) {
        ids.forEach((id) => {
            const target = state.list.find(i => i.id === id);
            if (target) {
                state.list.splice(state.list.indexOf(target), 1);
            }
        });
    },
    TOTAL(state, total) {
        state.total = total;
    },
    SET_QUEUED(state, count) {
        state.queued = count;
    },
    SET_TIMEOUTED(state, timeouted) {
        state.timeouted = timeouted;
    }
};
