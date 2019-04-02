module.exports = class Time {
    constructor() {
        this.start = new Date();
        this.times = {};
    }

    setTime(name) {
        if (name === 'total') {
            this.times[name] = ((new Date - this.start) / 1000).toFixed(4);
            return;
        }

        if (!this.times[name]) {
            this.times[name] = 0;
        }
        this.times[name] = (((new Date - this.start) / 1000) - this.sum()).toFixed(4);
    }

    getTime(name) {
        if(this.times[name]) {
            return String(this.times[name]).substring(0, 5);
        }
        return "0";
    }

    sum() {
        let total = 0;
        for(let key of Object.keys(this.times)) {
            total += parseFloat(this.times[key]);
        }
        return total;
    }
}
