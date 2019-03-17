'use strict';

const puppeteer = require('puppeteer');
const Analyzer = require('./Analyzer');
const Results = require('./Results');
const Pattern = require('./Pattern');

const tStart = new Date();

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
    waitUntil: 'networkidle0',
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
        page = await browser.newPage();
        await page.setUserAgent('Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36');
        await page.setRequestInterception(true);
        page.on('request', requestIntercept);
        await page.goto(website, gotoConf);
    } catch (e) {
        console.error(e)
        await browser.close();
    }

    try {
        matches = analyzer.analyze(page);
    } catch (e) {
        console.error(e)
        await browser.close();
    }

    const duration = (new Date() - tStart) / 1000;
    const res = new Results({ duration, matches });

    console.log(JSON.stringify(res));

    await browser.close();
})();
