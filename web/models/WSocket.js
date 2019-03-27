import Website from '@/models/Website';

const CONN_MSG = 'conn';
const COMM_MSG = 'comm';
const PING_MSG = 'ping';

class Msg {
    constructor({ message = '', userId = 0 } = {}) {
        this.message = message;
        this.userId = userId;
    }
}

class ConnMsg extends Msg {
    constructor(msg) {
        super(msg);
        this.type = CONN_MSG;
    }
}

class CommMsg extends Msg {
    constructor(msg) {
        super(msg);
        this.type = COMM_MSG;
    }
}

class PingMsg extends Msg {
    constructor(msg) {
        super(msg);
        this.type = PING_MSG;
    }
}

export default class WSocket {
    constructor({ url = 'ws://127.0.0.1:3033/api/ws', store = {}, notify = {} } = {}) {
        this.client = new WebSocket(url);
        this.client.onopen = this.init.bind(this);
        this.client.onmessage = this.onmessage.bind(this);
        this.store = store;
        this.notify = notify;
    }

    init() {
        this.send(new ConnMsg({ userId: this.userId() }));
        this.keepAlive();
    }

    onmessage(e) {
        const msg = JSON.parse(e.data);

        switch (msg.type) {
        case PING_MSG:
            console.log('pong recieved');
            break;
        case COMM_MSG:
            this.handleUserMsg(msg);
            break;
        default:
            console.warn('Unhandled ws message: ', e);
        }
    }

    handleUserMsg(msg) {
        if (msg.message.action == 'update:website') {
            const w = new Website(msg.message.website);
            this.store.commit('websites/UPDATE', w);
            this.notify.success({
                title: w.url,
                message: 'Inspect complete',
                position: 'bottom-right'
            });
        }
    }

    keepAlive() {
        setInterval(() => {
            this.send(new PingMsg({ userId: this.userId() }));
        }, 10000);
    }

    send(msg) {
        this.client.send(JSON.stringify(msg));
    }

    userId() {
        if (this.store.state.auth.user) {
            return this.store.state.auth.user.ID;
        }
        return 0;
    }
}
