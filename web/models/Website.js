import dayjs from 'dayjs';
import Tag from '@/models/Tag';

export default class Website {
    constructor({ id = null, url = '', inspectedAt = null, tags = [] } = {}) {
        this.ID = id;
        this.URL = url;
        this.Loading = false;
        this.Tags = tags.map(t => new Tag(t));
        this.InspectedAt = inspectedAt ? formatDate(inspectedAt) : null;
    }
}

const formatDate = (str) => {
    if (str) {
        return dayjs(str).format('YYYY-MM-DD HH:mm');
    }
    return '';
};
