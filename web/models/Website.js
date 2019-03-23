import dayjs from 'dayjs';

export default class Website {
    constructor({ ID = null, URL = '', InspectedAt = null }) {
        this.ID = ID;
        this.URL = URL;
        this.Loading = false;
        this.InspectedAt = InspectedAt ? formatDate(InspectedAt) : null;
    }
}

const formatDate = (str) => {
    if (str) {
        return dayjs(str).format('YYYY-MM-DD HH:mm');
    }
    return '';
};
