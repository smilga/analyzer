const toJSON = require('./utils').toJSON;

const JS_SOURCE = 'js_source';
const HTML = 'html';
const RESOURCE = 'resource';

class Match {
    constructor(id, value) {
        this.patternId = id;
        this.value = value;
    }
}

module.exports = class Analyzer {
    constructor(patterns) {
        this.patterns = patterns;
        this.resourceURLs = [];
    }

    analyzeResources() {
        const patterns = this.extractPatterns(RESOURCE);
        let matches = [];

        patterns.forEach(p => {
            matches = matches.concat(this.resourceURLMatch(p));
        });

        return matches;
    }

    async analyzeHTML(page) {
        const patterns = this.extractPatterns(HTML);
        let matches = await this.htmlMatch(page, patterns);

        return matches.map(m => new Match(m.id, JSON.stringify(m.el)));
    }

    async htmlMatch(page, patterns) {
        return await page.evaluate((toJSONFn, pat) => {
            const toJSON = new Function(' return (' + toJSONFn + ').apply(null, arguments)');
            const matches = [];

            pat.forEach(p => {
                let node = document.querySelector(p.value)
                if(node) {
                    matches.push({
                        el: toJSON.call(null, node),
                        id: p.id
                    });
                }
            });

            return matches;
        }, toJSON.toString(), patterns);
    }

    resourceURLMatch(pattern) {
        let matchedUrls = this.multiMatch(this.resourceURLs, pattern.value);

        return new Match(pattern.id, matchedUrls.join(","));
    }

    multiMatch(urls, pattern) {
        let matches = [];
        urls.forEach(url => {
            if(this.searchString(url, pattern)) {
                matches.push(url);
            }
        });
        return matches;
    }

    extractPatterns(type) {
        return this.patterns.filter(p => p.type === type);
    }

    searchString(str, pattern) {
        return new RegExp('^' + pattern.split(/\*+/).map(this.regExpEscape).join('.*') + '$').test(str);
    }

    regExpEscape(s) {
        return s.replace(/[|\\{}()[\]^$+*?.]/g, '\\$&');
    }
}
