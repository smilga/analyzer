module.exports = class Response {
    constructor({ duration, detected } = {}) {
        this.duration = duration;
        this.detectedServices = detected;
    }
}
