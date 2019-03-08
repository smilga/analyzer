const Pattern = require('./Pattern');

module.exports = class Service {
    constructor({ ID, Name, Patterns = [] } = {}) {
        this.id = ID;
        this.name = Name;
        this.patterns = Patterns.map(p => new Pattern(p));
    }
}
