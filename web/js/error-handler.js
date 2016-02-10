(function() {
    angular
        .module('ghca')
        .config(ConfigureErrorHandler);

    ConfigureErrorHandler.$inject = ['$httpProvider'];

    function ConfigureErrorHandler($httpProvider) {
        $httpProvider.interceptors.push(function($q, $rootScope, $log, $injector) {
             return {
                 'responseError': function(rejection) {
                     $log.debug(rejection);
                     $rootScope.errorDescription = rejection.data.error;
                     $injector.get('$uibModal').open({
                         templateUrl: '/static/bsod.html',
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
    }
})();
