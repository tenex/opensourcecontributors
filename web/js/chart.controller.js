(function() {
  angular
    .module('ghca')
    .controller('ChartController', ChartController);

  ChartController.$inject = ['$scope', '$state', 'Summary'];

  function ChartController($scope, $state, Summary) {
    var vm = this;

    vm.data = [];
    vm.options = {
      margin: { top: 20 },
      series: [
        {
          axis: "y",
          dataset: "contributions"
        }
      ],
      axes: {x: {key: "x", type: "date"}}
    };

    Summary.get({}, function(summary) {
      var ds = summary.dailySummary;
      vm.data = {
        contributions: ds
      }; //.map(function(x) { return x.count; });
    });
  }
})();
