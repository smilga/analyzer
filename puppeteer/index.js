const { Cluster } = require('puppeteer-cluster');
const Results = require('./Results');
const Time = require('./Time');
const config = require('./config');
const JobManager = require('./JobManager');
const ERR = require('./Errors');
const Scraper = require('./Scraper');

let cluster = {};
let manager = new JobManager;

(async () => {
    cluster = await Cluster.launch(config.cluster);

    start(cluster);

    cluster.on('taskerror', async (err, website) => {
        if(!err.message.includes(ERR.NOT_RESOLVED) && website.retry !== parseInt(process.env.RETRY)) {
            manager.storePending(website);
        } else {
            let results = new Results({
                websiteId: website.id,
                userId: website.userId,
                error: err
            });
            manager.storeResults(results);
        }
    });

    await cluster.task(async ({ page, data: website }) => {

        const scraper = new Scraper(page);

        const intercept = req => {
            if (req.resourceType() === 'script') {
                scraper.addScript(req.url());
                req.continue();
            } else if (req.resourceType() === 'image') {
                req.abort();
            } else {
                req.continue();
            }
        }

        const time = new Time()

        // Page settings
        await page.setDefaultTimeout(config.TIMEOUT);
        await page.setUserAgent(config.UA);
        await page.setRequestInterception(true);
        page.on('request', intercept);

        const response = await page.goto(website.url, config.goto);
        time.setTime('loaded');
        scraper.setHeaders(response.headers());

        await scraper.run();

        time.setTime('processed');
        time.setTime('total');

        let results = new Results({
            time: time,
            websiteId: website.id,
            userId: website.userId,
            ...scraper.results(),
        });

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

