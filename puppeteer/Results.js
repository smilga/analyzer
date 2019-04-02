const Time = require('./Time');

module.exports = class Results {
    constructor({ time = new Time, matches = [], websiteId, userId } = {}) {
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
            startedIn: this.time.getTime('started'),
            loadedIn: this.time.getTime('loaded'),
            resourceCheckIn: this.time.getTime('resourceCheck'),
            htmlCheckIn: this.time.getTime('htmlCheck'),
            totalIn: this.time.getTime('total')
        }
    }
}
