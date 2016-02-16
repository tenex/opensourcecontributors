(function() {
  angular
    .module('ghca')
    .config(ConfigureRoutes);

  ConfigureRoutes.$inject = ['$stateProvider', '$urlRouterProvider'];

  function ConfigureRoutes($stateProvider, $urlRouterProvider) {
    $urlRouterProvider.otherwise('/');
    $stateProvider
      .state('root', {
        url: '/',
        views: {
          search: {
            templateUrl: 'search-form.html',
            controller: 'SearchController',
            controllerAs: 'search'
          }
        }
      })
      .state('root.user', {
        url: 'user/{username}',
        abstract: true,
        views: {
          "@": {
            controller: 'UserController',
            controllerAs: 'userVm',
            templateUrl: 'user.common.html',
            resolve: {
              "username": function($stateParams) {
                return $stateParams.username;
              }
            }
          }
        }
      })
      .state('root.user.repositories', {
        url: '',
        templateUrl: 'user.repositories.html'
      })
      .state('root.user.events', {
        url: '/events/{page:int}',
        templateUrl: 'user.events.html'
      });
  }
})();
