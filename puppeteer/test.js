const { Cluster } = require('puppeteer-cluster');
const puppeteer = require('puppeteer');
const Redis = require('redis');
const client = Redis.createClient({ host: 'redis' });
const client2 = client.duplicate();
const Analyzer = require('./Analyzer');
const Results = require('./Results');
const Pattern = require('./Pattern');

const gotoConf = {
    timeout: 25000,
    waitUntil: 'networkidle2',
};

let cluster = {};

(async () => {
    cluster = await Cluster.launch({
        concurrency: Cluster.CONCURRENCY_CONTEXT,
        maxConcurrency: 4,
        puppeteerOptions: {
            args: ['--no-sandbox', '--disable-setuid-sandbox']
        }
    });

    startRedisPuller(cluster)

    cluster.on('taskerror', (err, data) => {
        console.log(`Error crawling ${data}: ${err.message}`);
    });

    await cluster.task(async ({ page, data: website }) => {
        const analyzer = new Analyzer(getPatterns());

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

        const tStart = new Date();
        const time = {
            started: "0",
            loaded: "0",
            resourceCheck: "0",
            htmlCheck: "0",
            total: "0"
        }

        let matches = [];

        time.started = ((new Date() - tStart) / 1000).toString();
        await page.setUserAgent('Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36');
        await page.setRequestInterception(true);
        page.on('request', requestIntercept);

        try {
            await page.goto(website.url, gotoConf);
        } catch(e) {
            let results = new Results({ time, websiteId: website.id, userId: website.userId });
            client2.lpush(['inspect:results', JSON.stringify(results)]);
            return;
        }

        time.loaded = (((new Date() - tStart) / 1000) - time.started).toString();

        matches = matches.concat(analyzer.analyzeSystem());

        matches = matches.concat(analyzer.analyzeResources());
        time.resourceCheck = (((new Date() - tStart) / 1000) - time.loaded - time.started).toFixed(3);

        matches = matches.concat(await analyzer.analyzeHTML(page));
        time.htmlCheck = (((new Date() - tStart) / 1000) - time.loaded - time.resourceCheck - time.started).toFixed(3);

        time.total = ((new Date() - tStart) / 1000).toString();

        let results = new Results({ time, matches, websiteId: website.id, userId: website.userId });

        console.log(JSON.stringify(results));

        client2.lpush(['inspect:results', JSON.stringify(results)]);
    });

    await cluster.idle();
    //await cluster.close();
})();

const startRedisPuller = cluster => {
    const brpop = () => {
        client.brpop("pending:websites", 5, function(err, value) {
            if (err) {
                console.error(err)
            }
            if(value) {
                let website = JSON.parse(value[1]);
                cluster.queue(website);
            }
            setTimeout( brpop, 1000 );
        });
    };
    brpop();
}

const getPatterns = () => {
    let patterns = [];
    client2.hgetall('inspect:patterns', function(err, obj){
        patterns = Object.values(obj).map(v => JSON.parse(v));
    });
    return patterns;
}
