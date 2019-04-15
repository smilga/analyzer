import Feature from '@/models/Feature';

export default class Service {
    constructor({ id = null, name = '', features = [] } = {}) {
        this.id = id;
        this.name = name;
        this.features = !features ? [] : features.map(f => new Feature(f));
    }
}
