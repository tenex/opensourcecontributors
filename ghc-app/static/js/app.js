(function() {
    var app = angular.module('ghca', [
        'ghcaServices',
        'angularMoment',
        'truncate',
        'ui.bootstrap'
    ]);

    app.config(['$httpProvider', function($httpProvider) {
        $httpProvider.interceptors.push(function($q, $rootScope, $log, $injector) {
             return {
                 'responseError': function(rejection) {
                     $log.debug(rejection);
                     $rootScope.errorDescription = rejection.data.error;
                     $injector.get('$uibModal').open({
                         templateUrl: 'bsod.html',
                         controller: 'BsodInstanceCtrl',
                         keyboard: true,
                         windowClass: 'bsod',
                         size: 'lg',
                         resolve: {
                             errorDescription: function() {
                                 return rejection.data.error;
                             }
                         }
                     });
                     return $q.reject(rejection);
                 }
             };
         });
    }]);

    app.controller('BsodInstanceCtrl', function($scope, $uibModalInstance, errorDescription) {
        $scope.errorDescription = errorDescription;
        $scope.ok = function() {
            $uibModalInstance.close();
        };
    });

    app.controller(
        "StatisticsController",
        ["$scope", "$log", "moment", "Statistics",
         function($scope, $log, moment, Statistics) {
             $scope.stats = Statistics.get(
                 {}, function(statsData) {
                     $scope.retrieved = true;
                 }
             );
         }]);

    app.controller("UserController", ["$scope", "$rootScope", "$log", "moment", "User", "Event", function($scope, $rootScope, $log, moment, User, Event) {
        $rootScope.errorDescription = '';
            $scope.eventPageSize = 50; // constant

            $scope.tabs = {
                none: 0,
                repoList: 1,
                eventList: 2
            };

            $scope.currentTab = $scope.tabs.repoList;
            $scope.isCurrentTab = function(t) {
                return $scope.currentTab === t;
            };
            $scope.setCurrentTab = function(t) {
                $scope.currentTab = t;
                if (t == $scope.tabs.eventList) {
                    $scope.getGHEvents();
                }
            };

            $scope.initialize = function() {
                $scope.username = "";
                $scope.processedUsername = ""; // The data below is for...
                $scope.userUrl = "";
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
                $scope.setCurrentTab($scope.tabs.repoList);
                $scope.clearEvents();

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
    ]);

    app.directive("eventOcticon", function() {

        var octiconMap = {
            "GollumEvent": "book",
            "IssuesEvent": "issue-opened",
            "PushEvent": "repo-push",
            "CommitCommentEvent": "comment",
            "ReleaseEvent": "tag",
            "PublicEvent": "megaphone",
            "MemberEvent": "person",
            "IssueCommentEvent": "comment-discussion"
        };

        var eventDescriptionMap = {
            "GollumEvent": "Wiki",
            "IssuesEvent": "Issue",
            "PushEvent": "Push",
            "CommitCommentEvent": "Commit Comment",
            "ReleaseEvent": "Release",
            "PublicEvent": "Repository made public",
            "MemberEvent": "Membership change",
            "IssueCommentEvent": "Issue comment"
        };

        return {
            restrict: "A",
            require: "^ngModel",
            scope: {
                ngModel: '='
            },
            template: '',
            link: function(scope, element, attrs) {
                element.addClass("octicon");
                element.addClass("octicon-" + octiconMap[scope.ngModel]);
                element.attr("data-toggle", "tooltip");
                element.attr("data-placement", "left");
                element.attr("title", eventDescriptionMap[scope.ngModel]);
                $(element).tooltip();
            }
        };

    });

})();
