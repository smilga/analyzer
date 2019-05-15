const Time = require('./Time');

module.exports = class Results {
    constructor({
        time = new Time,
        websiteId,
        userId,
        scripts = [],
        links = [],
        headers = [],
        cookies = [],
        html = '',
        description = '',
        error = null
    } = {}) {
        this.time = time;
        this.websiteId = websiteId;
        this.userId = userId;
        this.scripts = scripts;
        this.links = links;
        this.headers = headers;
        this.cookies = cookies;
        this.html = html;
        this.description = description;
        this.error = error;
    }
    toJSON() {
        let o = Object.assign(this, {
            scripts: this.scripts.map(s => ({ value: s })),
            links: this.links.map(l => ({ value: l })),
            headers: this.headers.map(h => ({
                key: Object.entries(h)[0][0],
                value: Object.entries(h)[0][1],
            })),
            html: { value: this.html },
            loadedIn: this.time.getTime('loaded'),
            processedIn: this.time.getTime('processed'),
            totalIn: this.time.getTime('total'),
        });
        delete o.time;
        return o;
    }
}
