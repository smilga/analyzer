import dayjs from 'dayjs';
import Tag from '@/models/Tag';

export default class Website {
    constructor({ id = null, url = '', inspectedAt = null, tags = [] } = {}) {
        this.id = id;
        this.url = url;
        this.loading = false;
        this.tags = tags.map(t => new Tag(t));
        this.inspectedAt = inspectedAt ? formatDate(inspectedAt) : null;
    }
}

const formatDate = (str) => {
    if (str) {
        return dayjs(str).format('YYYY-MM-DD HH:mm');
    }
    return '';
};
