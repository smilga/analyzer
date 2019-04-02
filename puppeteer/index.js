const { Cluster } = require('puppeteer-cluster');
const Redis = require('redis');
const client = Redis.createClient({ host: 'redis' });
const client2 = client.duplicate();
const Analyzer = require('./Analyzer');
const Results = require('./Results');
const Pattern = require('./Pattern');
const Time = require('./Time');
const gotoConf = require('./config').gotoConf;
const UA = require('./config').UA;
const clusterConf = require('./config').clusterConf;

let cluster = {};

(async () => {
    cluster = await Cluster.launch(clusterConf);

    startRedisPuller(cluster)

    cluster.on('taskerror', (err, data) => {
        console.log(`Error crawling ${data}: ${err.message}`);
    });

    await cluster.task(async ({ page, data: website }) => {
        const patterns = await getPatterns();
        const analyzer = new Analyzer(patterns);

        const requestIntercept = req => {
            if (req.resourceType() === 'image') {
                analyzer.resourceURLs.push(req.url());
                req.abort();
            } else if (req.resourceType() === 'script' || req.resourceType() === 'xhr') {
                analyzer.resourceURLs.push(req.url());
                req.continue();
            } else {
                req.continue();
            }
        }

        const time = new Time()
        let matches = [];

        time.setTime('started');
        await page.setUserAgent(UA);
        await page.setRequestInterception(true);
        page.on('request', requestIntercept);

        try {
            await page.goto(website.url, gotoConf);
        } catch(e) {
            let results = new Results({
                time: time,
                matches,
                websiteId: website.id,
                userId: website.userId
            });
            client2.lpush(['inspect:results', JSON.stringify(results)]);
            return;
        }

        time.setTime('loaded');
        matches = matches.concat(analyzer.analyzeSystem());
        matches = matches.concat(analyzer.analyzeResources());
        time.setTime('resourceCheck');
        matches = matches.concat(await analyzer.analyzeHTML(page));
        time.setTime('htmlCheck');
        time.setTime('total');

        let results = new Results({
            time: time,
            matches,
            websiteId: website.id,
            userId: website.userId
        });

        console.log(JSON.stringify(results));
        client2.lpush(['inspect:results', JSON.stringify(results)]);
    });

    await cluster.idle();
    //await cluster.close();
})();

const startRedisPuller = cluster => {
    const brpop = () => {

        if(!hasAwailableWorkers(cluster)) {
            console.log('no free workers waiting')
            setTimeout( brpop, 500 );
            return;
        }

        client.brpop("pending:websites", 5, function(err, value) {
            if (err) {
                console.error(err)
            }
            if(value) {
                let website = JSON.parse(value[1]);
                cluster.queue(website);
            }
            setTimeout( brpop, 500 );
        });
    };
    brpop();
}

const getPatterns = () => {
    return new Promise((res, rej) => {
        client2.hgetall('inspect:patterns', function(err, obj){
            if(err) {
                rej(err);
            }
            patterns = Object.values(obj).map(v => JSON.parse(v));
            patterns = patterns.map(p => new Pattern(p))
            res(patterns)
        });
    })
}

const hasAwailableWorkers = (c) => {
    return c.options.maxConcurrency > c.workersBusy.length;
}
