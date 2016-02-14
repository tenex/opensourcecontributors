(function() {
  angular
    .module('ghca')
    .controller('SearchController', SearchController);

  SearchController.$inject = ['$scope', '$state'];

  function SearchController($scope, $state) {
    var vm = this;

    vm.username = '';
    vm.doSearch = doSearch;

    //////////

    function doSearch() {
      $state.go('root.user.repositories', {
        username: vm.username
      });
    }
  }
})();
