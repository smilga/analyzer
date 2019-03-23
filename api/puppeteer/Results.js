module.exports = class Response {
    constructor({ time, matches } = {}) {
        this.time = time;
        this.matches = matches;
    }
    toJSON() {
        return {
            matches: this.matches,
            loadedIn: this.time.loaded,
            resourceCheckIn: this.time.resourceCheck,
            htmlCheckIn: this.time.htmlCheck,
            totalIn: this.time.total,
        }
    }
}
