<template>
  <div class="login">
    <h1>Login</h1>
    <el-form ref="form" :model="cred" label-width="120px">
      <el-form-item label="Email">
        <el-input v-model="cred.email" name="email" autocomplete="email" />
      </el-form-item>
      <el-form-item label="Password">
        <el-input v-model="cred.password" type="password" name="password" show-password />
      </el-form-item>
      <div class="action">
        <el-button @click="login">
          Login
        </el-button>
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
        };
    },
    methods: {
        login() {
            this.$axios.post('/api/login', this.cred)
                .then(res => this.auth(res.data))
                .catch((e) => {
                    this.$notify.error({
                        title: 'Error',
                        message: e.response.data.error,
                        position: 'bottom-right'
                    });
                });
        },
        async auth(token) {
            this.cred = { email: '', password: '' };
            this.$store.commit('auth/setToken', token);
            await this.$store.dispatch('auth/me');
            this.$router.push({ name: 'websites' });
        }
    }
};
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
