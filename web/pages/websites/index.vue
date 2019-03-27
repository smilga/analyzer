<template>
  <div class="websites">
    <div class="actions">
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
      <el-button class="rescan" icon="el-icon-refresh" @click="inspectAll">
        Inspect all
      </el-button>
    </div>
    <div class="filter">
      <el-select
        v-model="selected"
        class="filter-list"
        multiple
        filterable
        default-first-option
        placeholder="Filter"
      >
        <el-option
          v-for="item in filters"
          :key="item.name"
          :label="item.name"
          :value="item.id"
        />
      </el-select>
    </div>
    <div class="website-list">
      <el-table
        class="website-table"
        :data="websites"
        style="width: 100%"
      >
        <el-table-column
          prop="url"
          label="URL"
          width="200"
        />
        <el-table-column
          prop="inspectedAt"
          label="InspectedAt"
          width="300"
        />
        <el-table-column
          cell-class-name="service-column"
          prop="tags"
          label="Tags"
        >
          <template slot-scope="scope">
            <el-tag
              v-for="t in scope.row.tags"
              :key="scope.row.id + t.id"
              size="mini"
              height="25"
              class="tag-preview"
              type="info"
            >
              {{ t.value }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          prop=""
          width="70"
          label=""
        >
          <template slot-scope="scope">
            <nuxt-link :to="`/websites/${scope.row.id}/report`">
              <span class="icon-btn">
                <i class="el-icon-search" />
              </span>
            </nuxt-link>
            <span loading="true" class="icon-btn" @click="inspect(scope.row)">
              <i :class="[ scope.row.loading === true ? 'loading' : '', 'el-icon-refresh' ]" />
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
  </div>
</template>

<script>
import Website from '@/models/Website';
import Filter from '@/models/Filter';
import { mapState } from 'vuex';

export default {
    middleware: 'authenticated',
    data() {
        return {
            filters: [],
            selected: [],
            dialog: {
                addWebsite: false
            },
            websiteURL: ''
        };
    },
    watch: {
        selected: function () {
            this.fetch();
        }
    },
    computed: {
        ...mapState({
            websites: state => state.websites.list
        })
    },
    mounted() {
        this.fetch();
        this.$axios.get('/api/filters')
            .then(res => this.filters = res.data.map(p => new Filter(p)));
    },
    methods: {
        inspectAll() {
            this.websites.forEach((w) => {
                this.$store.commit('websites/SET_LOADING', { id: w.id, status: true });
            });
            this.$axios.get('/api/inspect/websites').then(console.log);
        },
        fetch() {
            this.$store.dispatch('websites/fetch', this.selected.join(','));
        },
        successImport(websites) {
            const ws = websites.map(w => new Website({ url: w.url }));
            this.$store.commit('websites/ADD', ws);
        },
        addWebsite() {
            const w = new Website({ url: this.websiteURL });
            this.$axios.post('/api/websites', w)
                .then((res) => {
                    w.id = res.data.id;
                    this.$store.commit('websites/ADD', [w]);
                    this.websiteURL = '';
                    this.dialog.addWebsite = false;
                });
        },
        inspect(website) {
            this.$store.commit('websites/SET_LOADING', { id: website.id, status: true });
            this.$axios.get(`/api/inspect/websites/${website.id}`)
                // .then(res => this.$store.commit('websites/UPDATE', new Website(res.data)))
                .catch((e) => {
                    this.$notify.error({
                        title: 'Error',
                        message: e.response.data.error,
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

.cell {
    display: flex !important;
}

.rescan {
    float: right;
}

.loading {
    -webkit-animation:spin 1s linear infinite;
    -moz-animation:spin 1s linear infinite;
    animation:spin 1s linear infinite;
}
.filter {
    margin: 50px 0 0 0;
    .el-select {
        width: 100%;
    }
}
.tag-preview {
    margin: 0 2px;
}
@-moz-keyframes spin { 100% { -moz-transform: rotate(360deg); } }
@-webkit-keyframes spin { 100% { -webkit-transform: rotate(360deg); } }
@keyframes spin { 100% { -webkit-transform: rotate(360deg); transform:rotate(360deg); } }
</style>
