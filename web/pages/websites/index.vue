<template>
  <div class="websites">
    <el-button icon="el-icon-circle-plus" @click="dialog.addWebsite = true">
      Add Website
    </el-button>
    <el-upload
      class="import"
      action="/api/websites/import"
      accept=".csv"
      :headers="{ Authorization: `Bearer ${$store.state.auth.token}` }"
      :on-success="successImport"
    >
      <el-button icon="el-icon-circle-plus">
        Import Websites
      </el-button>
      <div slot="tip" class="el-upload__tip">
        csv file with plain url list
      </div>
    </el-upload>
    <el-table
      class="website-table"
      :data="websites"
      style="width: 100%"
    >
      <el-table-column
        prop="URL"
        label="URL"
        width="200"
      />
      <el-table-column
        prop="SearchedAt"
        label="SearchedAt"
        width="100"
      />
      <el-table-column
        cell-class-name="service-column"
        prop="Services"
        label="Services"
      >
        <template slot-scope="scope">
          <img
            v-for="s in scope.row.Services"
            :key="scope.row.ID + s.ID"
            height="25"
            class="service-preview"
            :src="s.LogoURL"
            alt=""
          >
        </template>
      </el-table-column>
      <el-table-column
        prop=""
        width="50"
        label=""
      >
        <template slot-scope="scope">
          <span class="inspect" @click="inspect(scope.row)">
            <i :class="[ scope.row.Loading === true ? 'loading' : '', 'el-icon-refresh' ]" />
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
      <el-input v-model="websiteURL" placeholder="https://iprof.lv" />
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
    data() {
        return {
            websites: [],
            dialog: {
                addWebsite: false
            },
            websiteURL: ''
        };
    },
    asyncData({ app }) {
        return app.$axios.get('/api/websites')
            .then(res => ({ websites: res.data.map(w => new Website(w)) }));
    },
    methods: {
        successImport(websites) {
            this.websites = this.websites.concat(
                websites.map(w => new Website({ URL: w.URL }))
            );
        },
        addWebsite() {
            const w = new Website({ URL: this.websiteURL });
            this.$axios.post('/api/websites', w)
                .then((res) => {
                    w.ID = res.data.ID;
                    this.websites.push(w);
                    this.websiteURL = '';
                    this.dialog.addWebsite = false;
                });
        },
        inspect(website) {
            website.Loading = true;
            this.$axios.get(`/api/inspect/websites/${website.ID}`)
                .then((res) => {
                    website.SearchedAt = res.data.SearchedAt;
                    website.Services = res.data.Services;
                    website.Loading = false;
                })
                .catch((e) => {
                    this.$notify.error({
                        title: 'Error',
                        message: e.message,
                        position: 'bottom-right'
                    });
                });
        }
    }
};
</script>

<style lang="scss">
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

.cell {
    display: flex !important;
}

.loading {
    -webkit-animation:spin 1s linear infinite;
    -moz-animation:spin 1s linear infinite;
    animation:spin 1s linear infinite;
}
@-moz-keyframes spin { 100% { -moz-transform: rotate(360deg); } }
@-webkit-keyframes spin { 100% { -webkit-transform: rotate(360deg); } }
@keyframes spin { 100% { -webkit-transform: rotate(360deg); transform:rotate(360deg); } }
</style>
