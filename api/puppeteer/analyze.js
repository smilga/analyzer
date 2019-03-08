'use strict';

const puppeteer = require('puppeteer');
const Service = require('./Service');
const Analyzer = require('./Analyzer');

const args = process.argv.slice(2);
const website = args[0];
const service = new Service(JSON.parse(args[1]));

const analyzer = new Analyzer(service);

const connConfig = {
    ignoreHTTPSErrors: true,
    browserWSEndpoint: 'ws://browserless:3000',
};

const gotoConf = {
    timeout: 25000,
    waitUntil: 'networkidle0',
};

const requestIntercept = req => {
    if (req.resourceType() === 'script' || req.resourceType() === 'xhr') {
        console.log(req.url())
        //analyzer.resourceURLs.push(req.url());
    }
    req.continue();
}

(async() => {
    let browser = null;
    let page = null;
    let results = null;

    try {
        browser = await puppeteer.connect(connConfig);
        page = await browser.newPage();
        await page.setUserAgent('Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36');
        await page.setRequestInterception(true);
        page.on('request', requestIntercept);

        await page.goto(website, gotoConf);
        await page.waitForRequest('https://cdn.mxapis.com/mt.js?v=af3h');
    } catch (e) {
        console.error(e)
        await browser.close();
    }

    try {
        results = analyzer.analyze(page);
    } catch (e) {
        console.error(e)
        await browser.close();
    }

    console.log(results)

    await browser.close();
})();
