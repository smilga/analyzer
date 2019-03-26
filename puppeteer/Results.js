module.exports = class Response {
    constructor({ time = {}, matches = [], websiteId, userId } = {}) {
        this.time = time;
        this.matches = matches;
        this.websiteId = websiteId;
        this.userId = userId;
    }
    toJSON() {
        return {
            websiteId: this.websiteId,
            userId: this.userId,
            matches: this.matches,
            startedIn: this.time.started.substring(0, 4) || "0",
            loadedIn: this.time.loaded.substring(0, 4) || "0",
            resourceCheckIn: this.time.resourceCheck.substring(0, 4) || "0",
            htmlCheckIn: this.time.htmlCheck.substring(0, 4) || "0",
            totalIn: this.time.total.substring(0, 4) || "0",
        }
    }
}
