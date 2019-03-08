<template>
    <div class="create-edit-service">
        <el-form ref="form" :model="service" label-width="120px">
            <el-form-item label="Name">
                <el-input v-model="service.Name"></el-input>
            </el-form-item>
            <el-form-item label="Logo URL">
                <el-input v-model="service.LogoURL"></el-input>
            </el-form-item>
        </el-form>
        <div class="divider"></div>
        <div class="patterns">
            <div
                class="pattern"
                v-for="(p, i) in service.Patterns"
                :key="i"
            >
                <el-tooltip class="item" effect="light" content="Page uses service if this pattern match" placement="top-start">
                    <el-checkbox v-model="p.Mandatory">Mandatory</el-checkbox>
                </el-tooltip>
                <el-tag class="type" type="info">{{ p.Type }}</el-tag>
                <el-input :placeholder="placeholderByType(p.Type)" v-model="p.Value"></el-input>
                <el-button @click="deletePattern(p)" class="delete" icon="el-icon-minus" circle></el-button>
            </div>
        </div>
        <div class="create-patterns">
            <el-dropdown>
                <el-button>
                    Create new pattern
                    <i class="el-icon-arrow-down el-icon--right"></i>
                </el-button>
                <el-dropdown-menu slot="dropdown">
                    <el-dropdown-item
                        v-for="(t, i) in types"
                        :key="i"
                        @click.native="createPattern(t)"
                        >
                        {{ translate(t) }}
                    </el-dropdown-item>
                </el-dropdown-menu>
            </el-dropdown>
            <el-button @click="$emit('save')" icon="el-icon-check">Save</el-button>
        </div>
    </div>
</template>

<script>
import Pattern, { TYPE } from '@/models/Pattern';

export default {
    props: {
        service: {
            required: true,
            type: Object
        }
    },
    data() {
        return {
            types: TYPE,
        }
    },
    methods: {
        translate(type) {
            switch (type) {
                case TYPE.RESOURCE:
                    return 'Page loads resource from url';
                case TYPE.JS_SOURCE:
                    return 'Page scripts contains';
                case TYPE.HTML:
                    return 'Page has HTML tag';
                default:
                    return 'Err no name';
            }
        },
        placeholderByType(type) {
            switch (type) {
                case TYPE.RESOURCE:
                    return '*mt.js (accepts wildcard)';
                case TYPE.JS_SOURCE:
                    return '*jQuery* (accepts wildcard)';
                case TYPE.HTML:
                    return 'div.opt-in or div#someID (use CSS selectors)';
                default:
                    return 'Err no name';
            }
        },
        createPattern(type) {
            this.service.Patterns.push(new Pattern({
                ID: null,
                Type: type
            }));
        },
        deletePattern(p) {
            const target = this.service.Patterns.indexOf(p);
            this.service.Patterns.splice(target, 1);
        },
    }
}
</script>

<style lang="scss" scoped>
@import '@/assets/scss/_variables.scss';

.create-edit-service {
    max-width: 600px;
    margin: auto;
}
.create-patterns {
    justify-content: flex-end;
    display: flex;
    margin: 10px 0 0 0;
    .el-dropdown {
        margin-right: 10px;
    }
}
.pattern {
    margin-bottom: 11px;
    display: flex;
    align-items: center;
    justify-content: center;
    .delete {
        margin-left: 10px;
    }
    .type {
        margin-right: 10px;
    }
}
</style>
