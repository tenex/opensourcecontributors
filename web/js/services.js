(function() {
  var app = angular.module('ghca.services', [
    require('angular-resource')
  ]);
  app.factory('User', ['$resource', function($resource) {
    return $resource('/api/user/:username');
  }]);
  app.factory('Event', ['$resource', function($resource) {
    return $resource('/api/user/:username/events/:page');
  }]);
  app.factory('Statistics', ['$resource', function($resource) {
    return $resource('/api/stats');
  }]);
  app.factory('Summary', ['$resource', function($resource) {
    return $resource('/api/summaries');
  }]);
})();
