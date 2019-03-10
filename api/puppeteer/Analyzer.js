const JS_SOURCE = 'js_source';
const HTML = 'html';
const RESOURCE = 'resource';

class Match {
    constructor(id, value) {
        this.patternId = id;
        this.value = value;
    }
}

class DetectedService {
    constructor(id, match) {
        this.serviceId = id;
        this.match = match;
    }
}

// TODO check for Mandatory checkbox
module.exports = class Analyzer {
    constructor(services) {
        this.services = services;
        this.resourceURLs = [];
        this.detectedServices = [];
    }

    analyze(page) {
        this.services.forEach(s => this.scanForService(s));
        return this.detectedServices;
    }

    scanForService(service) {
        const patterns = this.extractPatterns(service, RESOURCE);
        let matched = [];

        patterns.forEach(pattern => {
            matched = matched.concat(this.resourceURLMatch(pattern));
        });

        if(matched.length > 0) {
            this.detectedServices.push(
                new DetectedService(service.id, matched)
            );
        }
    }

    resourceURLMatch(pattern) {
        const discovered = [];
        this.resourceURLs.forEach(url => {
            if(this.searchString(url, pattern.value)) {
                discovered.push(new Match(pattern.id, url));
            }
        });
        return discovered;
    }

    searchString(str, pattern) {
        return new RegExp('^' + pattern.replace(/\*/g, '.*') + '$').test(str);
    };

    extractPatterns(service, type) {
        return service.patterns.filter(p => p.type = type);
    }
}
