const { Cluster } = require('puppeteer-cluster');
const Analyzer = require('./Analyzer');
const Results = require('./Results');
const Time = require('./Time');
const config = require('./config');
const JobManager = require('./JobManager');
const ERR = require('./Errors');

let cluster = {};
let manager = new JobManager;

(async () => {
    cluster = await Cluster.launch(config.cluster);

    start(cluster);

    cluster.on('taskerror', async (err, website) => {
        if(err.message.includes(ERR.TIMEOUT) && website.retry !== parseInt(process.env.RETRY)) {
            console.log("TIMEOUT ERR: ", website.url)
            manager.storeTimeouted(website);
        } else {
            console.log('OTHER ERROR: ', website.url)
            let patterns = await manager.getPatterns();
            let results = new Results({
                time: new Time,
                matches: (new Analyzer(patterns)).getErrorMatch(err.message),
                websiteId: website.id,
                userId: website.userId
            });
            manager.storeResults(results);
        }
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
        await page.setDefaultTimeout(config.TIMEOUT);
        await page.setUserAgent(config.UA);
        await page.setRequestInterception(true);
        page.on('request', requestIntercept);

        await page.goto(website.url, config.goto);

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

        //console.log(JSON.stringify(results));
        manager.storeResults(results);
    });

    await cluster.idle();
    //await cluster.close();
})();

const hasAwailableWorkers = (c) => {
    return c.options.maxConcurrency > c.workersBusy.length;
}

const start = (cluster) => {
    console.log('start')
    const updateQueue = async () => {
        if(!hasAwailableWorkers(cluster)) {
            setTimeout( updateQueue, config.UPDATE_QUEUE_TIMEOUT );
            return;
        }

        let website = await manager.nextWebsite();
        if(website) {
            cluster.queue(website);
            setTimeout( updateQueue, config.UPDATE_QUEUE_TIMEOUT );
        } else {
            setTimeout( updateQueue, 100 );
        }
    }

    updateQueue();
}

