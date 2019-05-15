const Redis = require('redis');
const Pattern = require('./Pattern');
const Lists = require('./Lists');

const READ_LIST = process.env.READ_LIST;
const PENDING_LIST         = 'pending:websites:user:'
const LISTS_LIST           = 'pending:lists'
const TIMEOUTED_LIST       = 'timeouted:websites:user:'
const TIMEOUTED_LISTS_LIST = 'timeouted:lists'
const RESULTS_LIST         = 'inspect:results'
const PATTERNS_LIST        = 'inspect:patterns'

module.exports = class JobManager {
    constructor() {
        this.client = Redis.createClient({ host: 'redis' });
        this.lists = new Lists;
        setInterval(this.updateLists.bind(this), 1000);
    }

    nextWebsite() {
        return new Promise((res, rej) => {
            if(!this.lists.hasNext()) {
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
        this.client.smembers(READ_LIST, (err, val) => {
            if(err) {
                console.log('Error updating lists')
                return;
            }
            this.lists.setLists(val);
        });
    }

    storeResults(results) {
        let res = JSON.stringify(results)
        this.client.lpush([RESULTS_LIST, res]);
    }

    storePending(website) {
        if(website.retry) {
            website.retry++;
        } else {
            website.retry = 1;
        }

        const list = PENDING_LIST + website.userId;
        this.client.sadd(LISTS_LIST, list);
        this.client.lpush([list, JSON.stringify(website)]);
    }

    getPatterns() {
        return new Promise((res, rej) => {
            this.client.hgetall(PATTERNS_LIST, function(err, obj){
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
