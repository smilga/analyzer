export default class Website {
    constructor({ ID = null, URL = '', Services = [] }) {
        this.ID = ID;
        this.URL = URL;
        this.Services = Services;
        this.Loading = false;
    }
}
