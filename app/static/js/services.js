angular
	.module('ghcaServices', ['ngResource'])
	.factory('User', ['$resource',
		function($resource) {
			return $resource('/user/:username', {}, {});
		}])
	.factory('Event', ['$resource',
		function($resource) {
			return $resource('/user/:username/events/:page', {}, {});
		}]);
