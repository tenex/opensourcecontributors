(function() {
    angular
        .module('ghca')
        .controller('SearchController', SearchController);

    SearchController.$inject = ['$scope', '$state', '$log'];

    function SearchController($scope, $state, $log) {
        var vm = this;

        vm.username = '';

        vm.doSearch = doSearch;

        //////////

        function doSearch() {
            $state.go('user.summary', {
                username: vm.username
            });
        }
    }
})();
