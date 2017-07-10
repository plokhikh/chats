const DIRECTION_IN = 0,
      DIRECTION_OUT = 1;

var app = angular.module('chat', []);

app.factory('$socket', function() {
    return new Talker;
});

//watch selected user
app.directive('listUserItem', function () {
    return function (scope, element, attrs) {
        scope.listUserItemClass = function() {
            if (scope.currentUserId === element.data().userId) {
                return 'bg-success';
            } else {
                return '';
            }
        }

        element.bind("click", function (event) {
            scope.setcurrentUserId(element.data().userId);
            scope.$apply();
        });
    };
});

//watch user printed message
app.directive('chatInput', function () {
    return function (scope, element, attrs) {
        element.bind("keyup", function (event) {
            if ((event.keyCode == 10 || event.keyCode == 13)
//            && !event.ctrlKey
            ) {
                var msg = element.val().trim("\n");
                scope.send(msg);
                element.val('');
            }
        });
    };
});

app.controller('AppController', ['$scope', '$socket', function($scope, $socket) {
    var dummyUser = {"id": null, "name": null, "status": "online", "messages": []};
    var dummyMessage = {"direction": null, "content": null};
    $scope.currentUserId;
    $socket.onopen(function() {
        $socket.command("get-users-online");
    });

    //all connected users
    $scope.users = [];

    $scope.setcurrentUserId = function(userId) {
        $scope.currentUserId = userId;
    }

    $scope.getMessages = function(userId) {
        var messages = []
        $scope.users.forEach(function(user) {
            if (userId == user.id) {
                messages = user.messages;
            }
        });
        return messages;
    }

    $scope.getClassByDirection = function(direction) {
        return (direction == DIRECTION_IN && "pull-right") || (direction == DIRECTION_OUT && "pull-left");
    }

    function valueUsersOnline(data) {
        data.data.forEach(function(source) {
            pushToOnlineUsers(source.guid, source.name);
        });
    }

    /** функция, выполняющаяся если прилетело событие message-sended */
    function eventMessageSended(data) {
        $scope.users.forEach(function(user) {
            if (user.id == data.source.guid) {
                pushToMessages(user.id, DIRECTION_IN, data.data.message);
            }
        });

        $scope.$apply();
    }

    function eventUserEntered(data) {
        pushToOnlineUsers(data.source.guid, data.source.name);
    }

    function pushToOnlineUsers(guid, name) {
        var user = JSON.parse(JSON.stringify(dummyUser));

        user.id = guid;
        user.name = name;

        $scope.users.push(user);
        $scope.$apply();
    }

    function pushToMessages(userId, direction, message) {
        var msg = JSON.parse(JSON.stringify(dummyMessage));

        msg.direction = direction;
        msg.content = message;

        $scope.users.forEach(function(user) {
            if (userId == user.id) {
                user.messages.push(msg);
            }
        });

        $scope.$apply();
    }

    $scope.send = function(message) {
        pushToMessages($scope.currentUserId, DIRECTION_OUT, message);
        $socket.send($scope.currentUserId, message);
    };

    $socket.listen(function(string) {
        var message = JSON.parse(string);
        console.log(string);

        try {
            var prefix = ["event", "value"][message.type - 1];
            var func = prefix + snakeToCamel(message.name);
            (eval(func))(message);
        } catch(err) {
            console.error(err);
        }
    });
}]);

function snakeToCamel(s){
    var _s = s.replace(/(\-\w)/g, function(m){return m[1].toUpperCase();});
    return _s[0].toUpperCase() + _s.slice(1);
}