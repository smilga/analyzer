export default class Feature {
    constructor({
        id = null,
        value = '',
        createdAt = null
    } = {}) {
        this.id = id;
        this.value = value;
        this.createdAt = createdAt;
    }
}
