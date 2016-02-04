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
                templateUrl: '/static/user.common.html'
            })
            .state('user.repositories', {
                url: '/',
                controller: 'UserRepositoriesController',
                templateUrl: '/static/user.repositories.html'
            })
            .state('user.events', {
                url: '/events/{page}',
                controller: 'UserEventsController',
                templateUrl: '/static/user.events.html'
            });
            // .state('wtf', {

            // });
    }
})();
