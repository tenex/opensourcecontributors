(function() {
    angular
        .module('ghca')
        .controller("UserEventsController", UserEventsController);

    UserEventsController.$inject = [
        "$scope", "$rootScope", "$log", "User", "Event"
    ];

    function UserEventsController($scope, $rootScope, $log, User, Event) {
        $rootScope.errorDescription = '';
        $scope.eventPageSize = 50; // constant

        $scope.initialize = function() {
            $scope.username = "";
            $scope.processedUsername = ""; // The data below is for...
            $scope.eventCount = 0;
            $scope.repos = [];
            $scope.clearEvents();
        };

        $scope.clearEvents = function() {
            $scope.eventPages = {}; // a cache of sorts
            $scope.eventPageCount = 0;
            $scope.currentEventPage = 1;
            $scope.events = []; // the current page
            $scope.eventCount = 0;
        };

        $scope.initialize();

        // Have we retrieved the user's information (except a list of their events)?
        $scope.processed = false;
        $scope.processing = false;
        $scope.hasResults = false;
        $scope.loadingEvents = true;

        $scope.setCurrentEventsPage = function(i) {
            $scope.currentEventPage = i;
            $scope.getGHEvents();
        };

        $scope.eventPageChanged = function() {
            $scope.setCurrentEventsPage($scope.currentEventPage);
        };

        $scope.getGHEvents = function() {
            if ($scope.eventPages[$scope.currentEventPage]) {
                $scope.events = $scope.eventPages[$scope.currentEventPage];
                return;
            }

            $scope.loadingEvents = true;

            $scope.eventData = Event.get({
                username: $scope.processedUsername,
                page: $scope.currentEventPage
            }, function(eventData) {
                $scope.eventPages[$scope.currentEventPage] = eventData.events;
                $scope.events = eventData.events;
                $scope.loadingEvents = false;
            });
        };

        $scope.setUser = function() {
            $scope.processed = false;
            $scope.processing = true;

            $scope.user = User.get({
                username: $scope.username
            }, function(user) {
                $scope.processing = false;
                $scope.eventCount = user.eventCount;
                $scope.hasResults = user.eventCount ? true : false;
                $scope.eventPageCount = Math.ceil(
                    user.eventCount / $scope.eventPageSize);
                $scope.multipleEventPages = (
                    $scope.eventCount > $scope.eventPageSize);
                $scope.repos = user.repos;
                $scope.userUrl = "https://github.com/" + user.username;
                $scope.processedUsername = user.username;
                $scope.processed = true;
                $scope.processing = false;
            });
        };
    }
})();
