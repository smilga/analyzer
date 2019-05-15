module.exports = class Scraper {
    constructor(page) {
        this.page = page;
        this.scripts = [];
        this.links = [];
        this.headers = [];
        this.cookies = [];
        this.html = '';
        this.description = '';
    }

    async run() {
        this.cookies = await this.page.cookies();
        this.links = await this.page.$$eval('a', as => as.map(a => a.href));
        this.html = await this.page.content();
        try {
            this.description = await this.page.$eval('meta[name="description"]', element => element.textContent);
        } catch {
            this.description = '';
        }
    }

    addScript(url) {
        this.scripts.push(url);
    }

    setHeaders(headers = {}) {
        for(let key of Object.keys(headers)) {
            this.headers.push({
                [key]: headers[key]
            })
        }
    }

    results() {
        return {
            scripts: this.scripts,
            links: this.links,
            headers: this.headers,
            cookies: this.cookies,
            html: this.html,
            description: this.description
        }
    }
}
