<template>
    <div class="login">
        <h1>Login</h1>
        <el-form ref="form" :model="cred" label-width="120px">
            <el-form-item label="Email">
                <el-input name="email" autocomplete="email" v-model="cred.email"></el-input>
            </el-form-item>
            <el-form-item label="Password">
                <el-input type="password" name="password" v-model="cred.password"  show-password></el-input>
            </el-form-item>
            <div class="action">
                <el-button @click="login">Login</el-button>
            </div>
        </el-form>
    </div>
</template>

<script>
export default {
    data() {
        return {
            cred: {
                email: '',
                password: ''
            }
        }
    },
    methods: {
        login() {
            this.$axios.post('/api/login', this.cred)
                .then(res => this.auth(res.data));
        },
        async auth(token) {
            this.cred = { email: '', password: '' };
            this.$store.commit('auth/setToken', token)
            await this.$store.dispatch('auth/me');
            this.$router.push({ name: 'websites'});
        }
    }
}
</script>

<style lang="scss" scoped>
.login {
    max-width: 600px;
    margin: auto;
}

.action {
    display: flex;
    justify-content: flex-end;
}
</style>
