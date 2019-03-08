const Response = require('./Response');

const JS_SOURCE = 'js_source';
const HTML = 'html';
const RESOURCE = 'resource';

class Match {
    constructor(id, value) {
        this.id = id;
        this.value = value;
    }
}

module.exports = class Analyzer {
    constructor(service) {
        this.service = service;
        this.resourceURLs = [];
        this.matched = [];
    }

    analyze(page) {
        this.pattern(RESOURCE).forEach(p => this.resourceURLMatch(p));
        return this.matched;
    }

    resourceURLMatch(pattern) {
        this.resourceURLs.forEach(url => {
            if(this.searchString(url, pattern.value)) {
                this.matched = new Match(pattern.id, url);
            }
        });
    }

    searchString(str, pattern) {
        return new RegExp('^' + pattern.replace(/\*/g, '.*') + '$').test(str);
    };

    pattern(type) {
        return this.service.patterns.filter(p => p.type = type);
    }

    results() {
        return JSON.stringify(this.matched);
    }
}
