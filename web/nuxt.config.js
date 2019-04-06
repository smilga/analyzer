const pkg = require('./package');

module.exports = {
    env: {
        WEB_DOMAIN: process.env.WEB_DOMAIN
    },
    mode: 'universal',

    /*
     ** Headers of the page
     */
    head: {
        title: pkg.name,
        meta: [
            { charset: 'utf-8' },
            { name: 'viewport', content: 'width=device-width, initial-scale=1' },
            { hid: 'description', name: 'description', content: pkg.description }
        ],
        link: [
            { rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' }
        ]
    },

    loading: {
        color: '#909399',
        height: '2px'
    },

    eeroading: { color: '#fff' },

    css: [
        'element-ui/lib/theme-chalk/index.css'
    ],

    plugins: [
        '~/plugins/element-ui',
        '~/plugins/axios',
        { src: '~/plugins/toasted', ssr: false }
    ],

    modules: [
        '@nuxtjs/axios',
        '@nuxtjs/proxy'
    ],

    axios: {
        baseURL: `http://api:${process.env.API_PORT}`,
        proxy: true,
        withCredentials: true
    },

    proxy: {
        '/api': { target: `http://api:${process.env.API_PORT}`, secure: false, ws: true }
    },

    build: {
        transpile: [/^element-ui/]
    }
};
