import Tag from '@/models/Tag';

export const TYPE = {
    // JS_SOURCE: 'js_source',
    HTML: 'html',
    RESOURCE: 'resource'
};

export default class Pattern {
    constructor({
        ID = null,
        Type = TYPE.RESOURCE,
        Value = '',
        Description = '',
        Tags = [],
        CreatedAt = null
    } = {}) {
        this.id = ID;
        this.type = Type;
        this.value = Value;
        this.description = Description;
        this.tags = !Tags ? [] : Tags.map(t => new Tag(t));
        this.createdAt = CreatedAt;
    }
}
