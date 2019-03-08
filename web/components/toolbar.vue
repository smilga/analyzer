<template>
    <el-menu
        :default-active="$route.name"
        class="el-menu-top"
        background-color="white"
        text-color="#545c64"
        active-text-color="#000000">

        <template v-if="user">
            <nuxt-link to="/websites">
                <el-menu-item index="websites">
                    <i class="el-icon-menu"></i>
                    <span>Websites</span>
                </el-menu-item>
            </nuxt-link>

            <nuxt-link to="/services">
                <el-menu-item index="services">
                    <i class="el-icon-setting"></i>
                    <span>Services</span>
                </el-menu-item>
            </nuxt-link>
        </template>

        <div class="right">
            <template v-if="user">
                <el-dropdown class="user-dropdown">
                    <div class="el-dropdown-link">
                        {{ user.Email }}
                        <i class="el-icon-arrow-down el-icon--right"></i>
                    </div>
                    <el-dropdown-menu slot="dropdown">
                        <el-dropdown-item @click.native="logout">Logout</el-dropdown-item>
                    </el-dropdown-menu>
                </el-dropdown>
            </template>
            <template v-else>
                <nuxt-link to="/login">
                    <el-menu-item index="4">
                        <i class="el-icon-setting"></i>
                        <span>Login</span>
                    </el-menu-item>
                </nuxt-link>
            </template>
        </div>
    </el-menu>
</template>
<script>
import { mapState } from 'vuex';
export default {
    computed: {
        ...mapState({
            user: state => state.auth.user
        })
    },
    methods: {
        logout() {
            this.$axios.get('api/logout')
            .then(() => {
                this.$store.commit('auth/setUser', null);
                this.$store.commit('auth/setToken', null);
                this.$router.push({ name: 'login' })
            })
        }
    },
}
</script>

<style lang="scss">
@import '@/assets/scss/_variables.scss';

.el-menu-top {
    display: flex;
    @include shadow;
    li {
        font-size: 16px;
    }

}
.right {
    margin-left: auto;
}
.user-dropdown {
    line-height: 56px;
    margin-right: 20px;
    cursor: pointer;
    font-size: 16px;
}
</style>
