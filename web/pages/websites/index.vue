<template>
    <div class="websites">
        <el-button @click="dialog.addWebsite = true" icon="el-icon-circle-plus">Add Website</el-button>
        <el-upload
            class="import"
            action="/api/websites/import"
            accept=".csv"
            :headers="{ Authorization: `Bearer ${$store.state.auth.token}` }"
            :on-success="successImport"
            >
            <el-button icon="el-icon-circle-plus">Import Websites</el-button>
            <div slot="tip" class="el-upload__tip">csv file with plain url list</div>
        </el-upload>
        <el-table
            class="website-table"
            :data="websites"
            style="width: 100%">
            <el-table-column
                prop="URL"
                label="URL"
                width="200">
            </el-table-column>
            <el-table-column
                prop="SearchedAt"
                label="SearchedAt"
                width="100">
            </el-table-column>
            <el-table-column
                prop="Services"
                label="Services">
            </el-table-column>
            <el-table-column
                prop=""
                width="50"
                label="">
                <template slot-scope="scope">
                    <span class="inspect" @click="inspect(scope.row.ID)">
                        <i class="el-icon-refresh"></i>
                    </span>
                </template>
            </el-table-column>
        </el-table>

        <el-dialog
            title="Add Website"
            :visible="dialog.addWebsite"
            width="30%"
            >
            <span>Website URL</span>
            <el-input placeholder="https://iprof.lv" v-model="websiteURL"></el-input>
            <span slot="footer" class="dialog-footer">
                <el-button @click="dialog.addWebsite = false">Cancel</el-button>
                <el-button @click="addWebsite">Confirm</el-button>
            </span>
        </el-dialog>
    </div>
</template>

<script>
import Website from '@/models/Website';

export default {
    middleware: 'authenticated',
    asyncData ({ app }) {
        return app.$axios.get('/api/websites')
            .then(res => ({ websites: res.data }))
    },
    data() {
        return {
            websites: [],
            dialog: {
                addWebsite: false
            },
            websiteURL: '',
        }
    },
    methods: {
        successImport(websites) {
            this.websites = this.websites.concat(
                websites.map(w => new Website({URL: w.URL}))
            );
        },
        addWebsite() {
            const w = new Website({ URL: this.websiteURL });
            this.$axios.post('/api/websites', w)
                .then(res => {
                    this.websites.push(w);
                    this.websiteURL = '';
                    this.dialog.addWebsite = false;
                });
        },
        inspect(id) {
            this.$axios.get(`/api/inspect/websites/${id}`)
                .then(console.log)
        }
    }
}
</script>

<style lang="scss" scoped>
.website-table {
    margin-top: 20px;
}
.import {
    display: inline-flex;
}
.el-upload__tip {
    display: flex;
    align-items: center;
    margin-top: 0;
    margin-left: 5px;
}
.el-upload-list {
    margin-left: 30px;
}
.inspect {
    cursor: pointer;
}
</style>
