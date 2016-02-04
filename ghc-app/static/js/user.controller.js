(function() {
    angular
        .module('ghca')
        .controller("UserController",
                    UserController);

    UserController.$inject = [
        "$log", "User", "Event", "username"
    ];

    function UserController($log, User, Event, username) {
        var vm = this;

        // User stuff
        vm.username = username;
        vm.userUrl = "";
        vm.user = User.get({username: username}, function(user) {
            vm.userUrl = "https://github.com/" + user.username;
        });

        // Event stuff
        vm.loadingEvents = false;
        vm.eventPages = {};
        vm.eventData = {};
        vm.currentEventPage = 1;
        vm.eventPageSize = 50;
        vm.eventPageChanged = function() {
            fetchEventPage(vm.currentEventPage);
        };

        // Preload
        fetchEventPage(1);

        function fetchEventPage(pageNumber) {
            vm.loadingEvents = true;
            vm.eventData = Event.get({
                username: username,
                page: pageNumber
            }, function(events) {
                vm.eventPages[pageNumber] = events;
                vm.loadingEvents = false;
            });
        }
    }
})();
