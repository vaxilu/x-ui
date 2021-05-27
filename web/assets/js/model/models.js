class User {
    username = "";
    password = "";
}

class Msg {
    success = false;
    msg = "";
    obj = null;

    constructor(success, msg, obj) {
        if (success != null) {
            this.success = success;
        }
        if (msg != null) {
            this.msg = msg;
        }
        if (obj != null) {
            this.obj = obj;
        }
    }
}

class DBInbound {
    id = 0;
    userId = 0;
    up = 0;
    down = 0;
    remark = 0;
    enable = false;
    expiryTime = 0;
    listen = "";
    port = 0;
    protocol = "";
    settings = "";
    streamSettings = "";
    tag = "";
    sniffing = "";

    constructor(data) {
        if (data == null) {
            return;
        }
        ObjectUtil.cloneProps(this, data);
    }
}