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
    remark = "";
    enable = true;
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

    toInbound() {
        let settings = {};
        if (!ObjectUtil.isEmpty(this.settings)) {
            settings = JSON.parse(this.settings);
        }

        let streamSettings = {};
        if (!ObjectUtil.isEmpty(this.streamSettings)) {
            streamSettings = JSON.parse(this.streamSettings);
        }

        let sniffing = {};
        if (!ObjectUtil.isEmpty(this.sniffing)) {
            sniffing = JSON.parse(this.sniffing);
        }
        const config = {
            port: this.port,
            listen: this.listen,
            protocol: this.protocol,
            settings: settings,
            streamSettings: streamSettings,
            tag: this.tag,
            sniffing: sniffing,
        };
        return Inbound.fromJson(config);
    }

    hasLink() {
        switch (this.protocol) {
            case Protocols.VMESS:
            case Protocols.VLESS:
            case Protocols.TROJAN:
            case Protocols.SHADOWSOCKS:
                return true;
            default:
                return false;
        }
    }

    genLink(address="") {
        const inbound = this.toInbound();
        return inbound.genLink(address, this.remark);
    }
}

class AllSetting {
    webListen = "";
    webPort = 65432;
    webCertFile = "";
    webKeyFile = "";
    webBasePath = "/";

    xrayTemplateConfig = "";

    timeLocation = "Asia/Shanghai";

    constructor(data) {
        if (data == null) {
            return
        }
        ObjectUtil.cloneProps(this, data);
    }

    equals(other) {
        return ObjectUtil.equals(this, other);
    }
}