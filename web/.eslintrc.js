module.exports = {
    root: true,
    env: {
        browser: true,
        node: true
    },
    parserOptions: {
        parser: 'babel-eslint'
    },
    extends: [
        '@nuxtjs'
    ],
    rules: {
        "indent": ["error", 4],
        "semi": [2, "always"],
    }
}
