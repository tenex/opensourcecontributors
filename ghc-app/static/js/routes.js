(function() {
    angular
        .module('ghca')
        .config(ConfigureRoutes);

    ConfigureRoutes.$inject = ['$stateProvider', '$urlRouterProvider'];

    function ConfigureRoutes($stateProvider, $urlRouterProvider) {
        $urlRouterProvider.otherwise('/');
        $stateProvider
            .state('home', {
                url: '/'
            })
            .state('user', {
                url: '/user/{username}',
                abstract: true
            })
            .state('user.summary', {
                url: '/'
            })
            .state('user.events', {
                url: '/events/'
            });
    }
})();
