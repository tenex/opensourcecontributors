(function() {
    var app = angular.module('ghca', []);

    app.controller("UserController", ["$http","$log", function($http, $log) {
        this.user = "";
        this.userUrl = "";
        this.eventCount = 0;
        this.repos = [];
        // Have we retrieved the user's information (except all events)?
        this.processed = false;

        this.setUser = function() {
            var user = this;
            $http.get('/user/'+this.user, {}).success(function(data) {
                user = data;
                user.userUrl = "https://github.com/"+data.user;
                user.processed = true;
            });
        };
    }]);

    app.controller("EventController", ["$http","$log", function($http,$log) {
        this.event = {};
    }]);

})();
