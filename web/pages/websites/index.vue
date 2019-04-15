<template>
  <div class="websites" @keyup.enter="fetch">
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
          type="success"
          plain
          class="rescan"
          icon="el-icon-download"
          @click="exportWebsites"
        >
          Export filtered
        </el-button>
      </span>
    </div>
    <div class="filter" />

    <div class="table-actions">
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
      <el-input v-model="pagination.search" class="search-input" placeholder="Search" />
      <el-button
        type="primary"
        icon="el-icon-search"
        @click="fetch"
      />
    </div>
    <div class="website-list">
      <el-table
        v-loading="fetching"
        class="website-table"
        :data="websites"
        style="width: 100%"
      >
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
          width="90"
          label=""
        >
          <template slot-scope="scope">
            <nuxt-link :to="`/websites/${scope.row.id}/report`">
              <span class="icon-btn">
                <i class="el-icon-search" />
              </span>
            </nuxt-link>
            <span class="icon-btn" @click="inspect(scope.row)">
              <i class="el-icon-refresh" />
            </span>
            <nuxt-link class="icon-btn" to="#" @click.native="deleteTarget = scope.row">
              <i class="el-icon-delete" />
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
            fetching: false,
            pagination: {
                page: 1,
                limit: 10,
                search: ''
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
        inspect(website) {
            this.$axios.get(`/api/inspect/websites/${website.id}`)
                .then((res) => {
                    this.$notify.success({
                        title: 'Queued',
                        message: 'Website will be inspected',
                        position: 'bottom-right'
                    });
                })
                .catch((e) => {
                    this.$notify.error({
                        title: 'Error',
                        message: e.response.data.error,
                        position: 'bottom-right'
                    });
                });
        },
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
            this.$store.dispatch('websites/delete', this.deleteTarget.id)
                .then(() => this.deleteTarget = null);
        },
        fetch() {
            this.fetching = true;
            this.$store.dispatch('websites/fetch', {
                filters: this.selected.join(','),
                pagination: this.pagination
            }).then(this.fetching = false);
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
.search-input {
    max-width: 250px;
    margin: 0 5px;
    max-width: 250px;
}
.filter-list {
    flex-grow: 1;
}
</style>
