const COMMAND = 0,
        EVENT = 1
        VALUE = 2;

const SOURCE_TYPE_USER = 0;

class Talker {
    constructor() {
        var socket = new WebSocket("ws://127.0.0.1:8080");
        var that = this;
        this.onopenCallbacks = [];

        socket.onopen = function() {
            that.onopenCallbacks.forEach(function(callback) {
                callback();
            });
        };

        socket.onclose = function(event) {
            if (event.wasClean) {
                console.log('Соединение закрыто чисто');
            } else {
                console.error('Обрыв соединения'); // например, "убит" процесс сервера
            }
                console.error('Код: ' + event.code + ' причина: ' + event.reason);
        };

        socket.onerror = function(error) {
            console.error("Ошибка " + error.message);
        };

        this.onopen(function() {
            console.log("Соединение установлено.");
        });

        this.socket = socket;
    }

    onopen(callback) {
        this.onopenCallbacks.push(callback);
    }

    listen(callback) {
        this.socket.onmessage = function(event) {
            callback(event.data);
        }
    }

    command(name, data) {
        this.socket.send(JSON.stringify({"type": COMMAND, "ts": Date.now(), "name": name, "data": data}));
    }

    /** send message to socket */
    send(to, message) {
        this.command("send-message", {"to": {"type": SOURCE_TYPE_USER, "guid": to.toString()}, "message": message});
    }
}