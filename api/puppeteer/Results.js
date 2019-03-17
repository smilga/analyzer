module.exports = class Response {
    constructor({ duration, matches } = {}) {
        this.duration = duration;
        this.matches = matches;
    }
}
