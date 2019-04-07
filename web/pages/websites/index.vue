<template>
  <div class="websites">
    <div class="actions">
      <span>
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
      </span>
      <span>
        <el-button
          class="rescan"
          icon="el-icon-refresh"
          @click="inspectAll"
        >
          Inspect All
        </el-button>
        <el-button
          class="rescan"
          icon="el-icon-refresh"
          @click="inspectNew"
        >
          Inspect new
        </el-button>
        <el-button
          :disabled="selectedWebsites.length === 0"
          class="rescan"
          icon="el-icon-refresh"
          @click="inspectSelected"
        >
          <template v-if="selectedWebsites.length === 0">
            Select to inspect
          </template>
          <template v-else>
            Inspect {{ selectedWebsites.length }}
          </template>
        </el-button>
      </span>
    </div>
    <div class="filter">
      <el-select
        v-model="selected"
        class="filter-list"
        multiple
        filterable
        default-first-option
        placeholder="Filter websites"
      >
        <el-option
          v-for="item in filters"
          :key="item.name"
          :label="item.name"
          :value="item.id"
        />
      </el-select>
    </div>

    <div class="table-actions">
      <el-button
        type="danger"
        plain
        style="margin-right: 4px;"
        :disabled="selectedWebsites.length === 0"
        icon="el-icon-delete"
        @click="deleteTarget = selectedWebsites.slice()"
      >
        <template v-if="selectedWebsites.length === 0">
          Select to delete
        </template>
        <template v-else>
          Delete {{ selectedWebsites.length }}
        </template>
      </el-button>
      <el-button
        type="success"
        plain
        class="rescan"
        icon="el-icon-download"
        @click="exportWebsites"
      >
        Export filtered
      </el-button>
    </div>
    <div class="website-list">
      <el-table
        class="website-table"
        :data="websites"
        style="width: 100%"
        @selection-change="selectedWebsites = arguments[0]"
      >
        <el-table-column
          type="selection"
          width="35"
        />
        <el-table-column
          width="25"
        >
          <template slot-scope="scope">
            <span>
              <i v-if="scope.row.loading" class="el-icon-loading" />
            </span>
          </template>
        </el-table-column>
        <el-table-column
          prop="url"
          label="URL"
          width="250"
        />
        <el-table-column
          prop="inspectedAt"
          label="InspectedAt"
          width="150"
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
          width="40"
          label=""
        >
          <template slot-scope="scope">
            <nuxt-link :to="`/websites/${scope.row.id}/report`">
              <span class="icon-btn">
                <i class="el-icon-search" />
              </span>
            </nuxt-link>
          </template>
        </el-table-column>
      </el-table>
      <div class="pagination">
        <el-pagination
          layout="prev, pager, next, total"
          :total="total"
          :current-page.sync="pagination.page"
          :page-size="pagination.limit"
          @current-change="fetch"
        />
      </div>

      <el-dialog
        title="Delete confirmation"
        :visible="deleteTarget !== null"
        width="30%"
      >
        <span>Are you sure want to delete?</span>
        <span slot="footer" class="dialog-footer">
          <el-button @click="deleteTarget = null">Cancel</el-button>
          <el-button type="error" @click="deleteSelected">Confirm</el-button>
        </span>
      </el-dialog>

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
            pagination: {
                page: 1,
                limit: 10
            },
            selectedWebsites: [],
            filters: [],
            selected: [],
            dialog: {
                addWebsite: false
            },
            websiteURL: '',
            deleteTarget: null
        };
    },
    watch: {
        selected: function () {
            // this.pagination.page = 1;
            this.fetch();
        }
    },
    computed: {
        ...mapState({
            websites: state => state.websites.list,
            total: state => state.websites.total
        })
    },
    mounted() {
        this.fetch();
        this.$axios.get('/api/filters')
            .then(res => this.filters = res.data.map(p => new Filter(p)));
    },
    methods: {
        exportWebsites() {
            return this.$axios.post(`/api/inspect/export?f=${this.selected.join(',')}`)
                .then((res) => {
                    const csv = 'data:text/csv;charset=utf-8,' + res.data;
                    const encodedUri = encodeURI(csv);
                    const link = document.createElement('a');
                    link.setAttribute('href', encodedUri);
                    link.setAttribute('download', 'results.csv');
                    document.body.appendChild(link); // Required for FF
                    link.click(); // This will download the data file named "my_data.csv".
                });
        },
        inspectAll() {
            return this.$axios.get('api/inspect/websites')
                .then((res) => {
                    this.$notify.success({
                        title: 'Start inspect',
                        message: `${res.data.count} websites will be inspected`,
                        position: 'bottom-right'
                    });
                });
        },
        inspectNew() {
            return this.$axios.get('api/inspect/new')
                .then((res) => {
                    this.$notify.success({
                        title: 'Start inspect',
                        message: `${res.data.count} websites will be inspected`,
                        position: 'bottom-right'
                    });
                });
        },
        deleteSelected() {
            const ids = this.selectedWebsites.map(w => w.id);
            this.$store.dispatch('websites/delete', ids)
                .then(() => {
                    this.deleteTarget = null;
                    this.selectedWebsites = [];
                });
        },
        inspectSelected() {
            this.selectedWebsites.forEach((w) => {
                this.$store.commit('websites/SET_LOADING', { id: w.id, status: true });
            });
            const ids = this.selectedWebsites.map(w => w.id);
            this.$axios.post('/api/inspect/websites', ids).then(() => {
            });
        },
        fetch() {
            this.$store.dispatch('websites/fetch', {
                filters: this.selected.join(','),
                pagination: this.pagination
            });
        },
        successImport(websites) {
            this.fetch();
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
.actions {
    display: flex;
    justify-content: space-between;
}
.pagination {
    text-align: center;
    margin: 20px;
}
.table-actions {
    margin-top: 10px;
    display: flex;
    justify-content: space-between;
}

span.el-pagination__total{
    display: block !important;
}
@-moz-keyframes spin { 100% { -moz-transform: rotate(360deg); } }
@-webkit-keyframes spin { 100% { -webkit-transform: rotate(360deg); } }
@keyframes spin { 100% { -webkit-transform: rotate(360deg); transform:rotate(360deg); } }
</style>
