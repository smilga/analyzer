export const TYPE = {
    JS_SOURCE: 'js_source',
    HTML: 'html',
    RESOURCE: 'resource',
}

export default class Pattern {
    constructor({ID = null, Mandatory = false, Type, Value = ''} = {}) {
        this.ID = ID;
        this.Mandatory = Mandatory
        this.Type = Type;
        this.Value = Value;
    }
}
