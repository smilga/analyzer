import cookieparser from 'cookieparser';

export const actions = {
    async nuxtServerInit ({commit, dispatch}, {req}) {

        let token = null;

        if (req.headers.cookie) {
            var parsed = cookieparser.parse(req.headers.cookie);

            if(parsed['access_token']) {
                token = parsed['access_token'];
            }

            commit('auth/setToken', token);
        }

        if(token && token.length > 0) {
            await dispatch('auth/me');
        }
    }
};
