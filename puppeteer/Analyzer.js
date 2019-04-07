const toJSON = require('./utils').toJSON;

const JS_SOURCE = 'js_source';
const HTML = 'html';
const RESOURCE = 'resource';
const SYSTEM = 'system';

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
            let match = this.resourceURLMatch(p)
            if(match) {
                matches.push(match);
            }
        });

        return matches;
    }

    analyzeSystem() {
        const patterns = this.extractPatterns(SYSTEM);
        let matches = [];


        patterns.forEach(p => {
            let match = this.systemMatch(p);
            if (match) {
                matches.push(new Match(p.id, match.value));
            }
        });

        return matches;
    }

    systemMatch(pattern) {
        if(pattern.value === 'isAlive') {
            return { value: "true" }
        }
    }

    getErrorMatch(err) {
        const patterns = this.extractPatterns(SYSTEM);
        let matches = [];
        patterns.forEach(p => {
            if(p.value === 'hasError') {
                matches.push(new Match(p.id, err))
            }
        })
        return matches;
    }

    async analyzeHTML(page) {
        const patterns = this.extractPatterns(HTML);
        let matches = await this.htmlMatch(page, patterns);

        return matches.map(m => new Match(m.id, JSON.stringify(m.el)));
    }

    async htmlMatch(page, patterns) {
        page.on('console', consoleObj => console.log(consoleObj.text()));

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

        if(matchedUrls.length > 0) {
            return new Match(pattern.id, matchedUrls[0]);
        }
        return null;
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
