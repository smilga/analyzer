const { Cluster } = require('puppeteer-cluster');

const goto = {
    waitUntil: 'networkidle2',
};

const UPDATE_QUEUE_TIMEOUT = 100;
const WORKERS = process.env.CONCURRENCY;
const TIMEOUT = process.env.TIMEOUT;

const UA = 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36';

const cluster = {
    concurrency: Cluster.CONCURRENCY_CONTEXT,
    maxConcurrency: WORKERS,
    puppeteerOptions: {
        args: ['--no-sandbox', '--disable-setuid-sandbox', '--disable-http2'],
        timeout: 30000, // timeout to start browser instance
    },
    timeout: TIMEOUT, // Specify a timeout for all task
}

const read = () => {
    console.log("READING ENVIROMENTS")
    console.log(process.env)
}
read();

module.exports = {
    goto,
    UA,
    cluster,
    WORKERS,
    UPDATE_QUEUE_TIMEOUT,
    TIMEOUT
}
