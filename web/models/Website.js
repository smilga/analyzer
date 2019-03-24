import dayjs from 'dayjs';
import Tag from '@/models/Tag';

export default class Website {
    constructor({ ID = null, URL = '', InspectedAt = null, Tags = [] } = {}) {
        this.ID = ID;
        this.URL = URL;
        this.Loading = false;
        this.Tags = Tags.map(t => new Tag(t));
        this.InspectedAt = InspectedAt ? formatDate(InspectedAt) : null;
    }
}

const formatDate = (str) => {
    if (str) {
        return dayjs(str).format('YYYY-MM-DD HH:mm');
    }
    return '';
};
