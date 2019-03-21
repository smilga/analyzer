module.exports = class Response {
    constructor({ time, matches } = {}) {
        this.time = time;
        this.matches = matches;
    }
    toJSON() {
        return {
            matches: this.matches,
            loadedIn: this.loaded,
            resourceCheckIn: this.resourceCheck,
            htmlCheckIn: this.htmlCheck,
            totalIn: this.totalIn,
        }
    }
}
