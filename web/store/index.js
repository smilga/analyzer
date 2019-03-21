import cookieparser from 'cookieparser';

export const actions = {
    async nuxtServerInit({ commit, dispatch }, { req }) {
        let token = null;

        if (req.headers.cookie) {
            const parsed = cookieparser.parse(req.headers.cookie);

            if (parsed.access_token) {
                token = parsed.access_token;
            }

            commit('auth/setToken', token);
        }

        if (token && token.length > 0) {
            try {
                await dispatch('auth/me');
            } catch (e) {
                console.error('Error loading me: ', e);
            }
        }
    }
};
