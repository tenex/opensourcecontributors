(function() {
  var app = angular.module('ghcaServices', ['ngResource']);
  app.factory('User', ['$resource', function($resource) {
    return $resource('/user/:username');
  }]);
  app.factory('Event', ['$resource', function($resource) {
    return $resource('/user/:username/events/:page');
  }]);
  app.factory('Statistics', ['$resource', function($resource) {
    return $resource('/stats');
  }]);
})();
