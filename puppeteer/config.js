const { Cluster } = require('puppeteer-cluster');

const gotoConf = {
    waitUntil: 'networkidle2',
};

const WORKERS = 16;
const UPDATE_QUEUE_TIMEOUT = 100;
const PAGE_LOAD_TIMEOUT = 60000;

const UA = 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36';

const clusterConf = {
    concurrency: Cluster.CONCURRENCY_CONTEXT,
    maxConcurrency: WORKERS,
    puppeteerOptions: {
        args: ['--no-sandbox', '--disable-setuid-sandbox', '--disable-http2'],
        timeout: 60000,
    },
    timeout: 60000,
    retryLimit: 3,
}

module.exports = {
    gotoConf,
    UA,
    clusterConf,
    WORKERS,
    UPDATE_QUEUE_TIMEOUT,
    PAGE_LOAD_TIMEOUT
}
