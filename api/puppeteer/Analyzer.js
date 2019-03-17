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

    analyze(page) {
        const resourcePatterns = this.extractPatterns(RESOURCE);

        let matches = [];
        resourcePatterns.forEach(p => {
            matches = matches.concat(this.resourceURLMatch(p));
        });

        // get other pattern types and analyze them

        return matches;
    }

    resourceURLMatch(pattern) {
        const matches = [];
        this.resourceURLs.forEach(url => {
            if(this.searchString(url, pattern.value)) {
                matches.push(new Match(pattern.id, url));
            }
        });
        return matches;
    }

    searchString(str, pattern) {
        // TODO use some library
        return new RegExp('^' + pattern.replace(/\*/g, '.*') + '$').test(str);
    };

    extractPatterns(service, type) {
        return this.patterns.filter(p => p.type = type);
    }
}
