const Redis = require('redis');
const Pattern = require('./Pattern');

const LISTS   = "pending:lists";
const RESULTS_LIST = "inspect:results";

class Lists {
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

module.exports = class JobManager {
    constructor() {
        this.client = Redis.createClient({ host: 'redis' });
        this.lists = new Lists;
        setInterval(this.updateLists.bind(this), 1000);
    }

    nextWebsite() {
        return new Promise((res, rej) => {
            if(!this.lists.hasNext()) {
                console.log('there is no lists', this.lists.lists)
                return res(null);
            }

            this.client.rpop(this.lists.next(), function(err, value) {
                if (err) {
                    return rej(err);
                }

                if(value) {
                    return res(JSON.parse(value));
                }
                return res(null);
            });
        });
    }

    updateLists() {
        this.client.smembers(LISTS, (err, val) => {
            if(err) {
                console.log('Error updating lists')
                return;
            }
            this.lists.setLists(val);
        });
    }

    storeResults(results) {
        this.client.lpush([RESULTS_LIST, JSON.stringify(results)]);
    }

    getPatterns() {
        return new Promise((res, rej) => {
            this.client.hgetall('inspect:patterns', function(err, obj){
                if(err) {
                    return rej(err);
                }
                let patterns = Object.values(obj).map(v => JSON.parse(v));
                patterns = patterns.map(p => new Pattern(p))
                return res(patterns)
            });
        })
    }
}
