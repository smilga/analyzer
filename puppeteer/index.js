const { Cluster } = require('puppeteer-cluster');
const Analyzer = require('./Analyzer');
const Results = require('./Results');
const Time = require('./Time');
const gotoConf = require('./config').gotoConf;
const UA = require('./config').UA;
const clusterConf = require('./config').clusterConf;
const UPDATE_QUEUE_TIMEOUT = require('./config').UPDATE_QUEUE_TIMEOUT;
const JobManager = require('./JobManager');

let cluster = {};
let manager = new JobManager;

(async () => {
    cluster = await Cluster.launch(clusterConf);

    start(cluster);

    cluster.on('taskerror', (err, data) => {
        console.log(`Error crawling ${data}: ${err.message}`);
    });

    await cluster.task(async ({ page, data: website }) => {
        const patterns = await manager.getPatterns();
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
            manager.storeResults(results);
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
        manager.storeResults(results);
    });

    await cluster.idle();
    //await cluster.close();
})();

const hasAwailableWorkers = (c) => {
    return c.options.maxConcurrency > c.workersBusy.length;
}

const start = (cluster) => {
    const updateQueue = async () => {
        if(!hasAwailableWorkers(cluster)) {
            setTimeout( updateQueue, UPDATE_QUEUE_TIMEOUT );
            return;
        }

        let website = await manager.nextWebsite();
        if(website) {
            cluster.queue(website);
            setTimeout( updateQueue, UPDATE_QUEUE_TIMEOUT );
        } else {
            setTimeout( updateQueue, 100 );
        }
    }

    updateQueue();
}

