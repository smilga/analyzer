export default class Website {
    constructor({ ID = null, URL = '', Services = [], InspectedAt = null }) {
        this.ID = ID;
        this.URL = URL;
        this.Services = Services;
        this.Loading = false;
        this.InspectedAt = InspectedAt;
    }
}
