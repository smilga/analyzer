export default class Service {
    constructor({ID = null, Name = '', LogoURL = '', Patterns = []} = {}) {
        this.ID = ID;
        this.Name = Name;
        this.LogoURL = LogoURL;
        this.Patterns = Patterns;
    }
}
