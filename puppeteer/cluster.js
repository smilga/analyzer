const { Cluster } = require('puppeteer-cluster');

const gotoConf = {
    timeout: 25000,
    waitUntil: 'networkidle2',
};

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


(async () => {
  const cluster = await Cluster.launch({
    concurrency: Cluster.CONCURRENCY_CONTEXT,
    maxConcurrency: 2,
  });

  await cluster.task(async ({ page, data: url }) => {
    let matches = [];

    try {
        time.started = ((new Date() - tStart) / 1000).toString();
        await page.setUserAgent('Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36');
        await page.setRequestInterception(true);
        page.on('request', requestIntercept);
        await page.goto(website, gotoConf);

        time.loaded = (((new Date() - tStart) / 1000) - time.started).toString();
    } catch (e) {
        if(e.message.includes('net::ERR_NAME_NOT_RESOLVED')) {
            console.log(JSON.stringify(new Results));
            return;
        }
    }

    matches = matches.concat(analyzer.analyzeSystem());

        matches = matches.concat(analyzer.analyzeResources());
        time.resourceCheck = (((new Date() - tStart) / 1000) - time.loaded - time.started).toFixed(3);

        matches = matches.concat(await analyzer.analyzeHTML(page));
        time.htmlCheck = (((new Date() - tStart) / 1000) - time.loaded - time.resourceCheck - time.started).toFixed(3);

    time.total = ((new Date() - tStart) / 1000).toString();

    console.log(JSON.stringify(new Results({ time, matches })));
  });

  cluster.queue('http://www.google.com/');
  cluster.queue('http://www.wikipedia.org/');
  // many more pages

  await cluster.idle();
  //await cluster.close();
})();

