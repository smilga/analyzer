export default class Filter {
    constructor({
        ID = null,
        Name = '',
        Description = '',
        Tags = []
    } = {}) {
        this.id = ID;
        this.name = Name;
        this.Description = Description;
        this.tags = Tags;
    }
}
