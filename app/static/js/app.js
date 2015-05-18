(function() {
    var app = angular.module('ghca', []);
    app.controller("EventController", function() {
        this.event = _event;
    });

    var _event = {
        _user_lower: "hut8",
        type: "PushEvent"
    };
})();
