export default class Tag {
    constructor({
        ID = null,
        Value = '',
        CreatedAt = null
    } = {}) {
        this.id = ID;
        this.value = Value;
        this.createdAt = CreatedAt;
    }
}
