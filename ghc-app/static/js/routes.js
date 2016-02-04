(function() {
    angular
        .module('ghca')
        .config(ConfigureRoutes);

    ConfigureRoutes.$inject = ['$stateProvider', '$urlRouterProvider'];

    function ConfigureRoutes($stateProvider, $urlRouterProvider) {
        $urlRouterProvider.otherwise('/');
        $stateProvider
            .state('home', {
                url: '/',
                template: '<div class="well"><h3>Search for a user above to begin</h3></div>'
            })
            .state('user', {
                url: '/user/{username}',
                abstract: true,
                controller: 'UserController',
                controllerAs: 'userVm',
                templateUrl: '/static/user.common.html',
                resolve: {
                    "username": function($stateParams) {
                        return $stateParams.username;
                    }
                }
            })
            .state('user.repositories', {
                url: '/',
                templateUrl: '/static/user.repositories.html'
            })
            .state('user.events', {
                url: '/events/{page:int}',
                controller: 'UserEventsController',
                templateUrl: '/static/user.events.html'
            });
            // .state('wtf', {

            // });
    }
})();
