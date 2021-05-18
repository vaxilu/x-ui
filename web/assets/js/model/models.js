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