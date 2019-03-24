'use strict';

const puppeteer = require('puppeteer');
const Analyzer = require('./Analyzer');
const Results = require('./Results');
const Pattern = require('./Pattern');

const tStart = new Date();

// TODO to class to handle times
const time = {
    started: null,
    loaded: null,
    resourceCheck: null,
    htmlCheck: null,
    total: null
}

const args = process.argv.slice(2);
const website = args[0];
const patterns = JSON.parse(args[1]).map(s => new Pattern(s));
const analyzer = new Analyzer(patterns);

const connConfig = {
    ignoreHTTPSErrors: true,
    browserWSEndpoint: 'ws://browserless:3000',
};

const gotoConf = {
    timeout: 25000,
    waitUntil: 'networkidle2',
};

const requestIntercept = req => {
    if (req.resourceType() === 'script' || req.resourceType() === 'xhr') {
        analyzer.resourceURLs.push(req.url());
    }
    req.continue();
}

(async() => {
    let browser = null;
    let page = null;
    let matches = [];

    try {
        browser = await puppeteer.connect(connConfig);
        time.started = ((new Date() - tStart) / 1000).toString();
        page = await browser.newPage();
        await page.setUserAgent('Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36');
        await page.setRequestInterception(true);
        page.on('request', requestIntercept);
        await page.goto(website, gotoConf);

        time.loaded = (((new Date() - tStart) / 1000) - time.started).toString();
    } catch (e) {
        console.error(e)
        await browser.close();
    }

    try {
        matches = analyzer.analyzeResources();
        time.resourceCheck = (((new Date() - tStart) / 1000) - time.loaded - time.started).toFixed(3);
    } catch (e) {
        console.error(e)
        await browser.close();
    }

    try {
        matches = matches.concat(await analyzer.analyzeHTML(page));
        time.htmlCheck = (((new Date() - tStart) / 1000) - time.loaded - time.resourceCheck - time.started).toFixed(3);
    } catch (e) {
        console.error(e);
        await browser.close();
    }

    time.total = ((new Date() - tStart) / 1000).toString();

    const res = new Results({ time, matches });

    console.log(JSON.stringify(res));

    await browser.close();
})();
