(function() {
  angular
    .module('ghca')
    .controller('ChartController', ChartController);

  ChartController.$inject = ['$scope', '$state', 'Summary'];

  function ChartController($scope, $state, Summary) {
    var vm = this;

    vm.data = [];
    vm.labels = [];
    vm.chartConfig = {
      xAxes: [{
        display: false
      }]
    };

    Summary.get({}, function(summary) {
      var ds = summary.dailySummary;
      vm.data = ds.map(function(x) { return x.count; });
      vm.labels = ds.map(function(x) { return x.date; });
    });
  }
})();
