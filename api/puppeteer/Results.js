module.exports = class Response {
    constructor({ time, matches } = {}) {
        this.time = time;
        this.matches = matches;
    }
    toJSON() {
        return {
            matches: this.matches,
            startedIn: this.time.started,
            loadedIn: this.time.loaded,
            resourceCheckIn: this.time.resourceCheck,
            htmlCheckIn: this.time.htmlCheck,
            totalIn: this.time.total,
        }
    }
}
