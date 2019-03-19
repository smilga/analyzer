module.exports = class Response {
    constructor({ time, matches } = {}) {
        this.time = time;
        this.matches = matches;
    }
}
