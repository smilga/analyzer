module.exports = class Lists {
    constructor() {
        this.lists = [];
        this.index = 0;
    }

    setLists(lists = []) {
        this.lists = lists.sort();
    }

    hasNext() {
        return this.lists.length > 0;
    }

    next() {
        if (this.index === this.lists.length) {
            this.index = 0;
        }

        return this.lists[this.index++];
    }
}
