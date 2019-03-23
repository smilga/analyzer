import Tag from '@/models/Tag';

export default class Filter {
    constructor({
        ID = null,
        Name = '',
        Description = '',
        Tags = []
    } = {}) {
        this.id = ID;
        this.name = Name;
        this.description = Description;
        this.tags = Tags.map(t => new Tag(t));
    }
}
