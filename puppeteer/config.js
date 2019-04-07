const { Cluster } = require('puppeteer-cluster');

const gotoConf = {
    timeout: 25000,
    waitUntil: 'networkidle2',
};

const WORKERS = 18;

const UA = 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36';

const clusterConf = {
    concurrency: Cluster.CONCURRENCY_CONTEXT,
    maxConcurrency: WORKERS,
    puppeteerOptions: {
        args: ['--no-sandbox', '--disable-setuid-sandbox']
    }
}

module.exports = {
    gotoConf,
    UA,
    clusterConf,
    WORKERS
}
