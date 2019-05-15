export default ({ app, store }) => {
    app.$axios.interceptors.request.use(function (config) {
        config.headers.Authorization = `Bearer ${store.state.auth.token}`;
        return config;
    }, null);
};
