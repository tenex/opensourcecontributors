(function() {
    angular
        .module('ghca')
        .controller("StatisticsController", StatisticsController);

    StatisticsController.$inject = ["$scope", "$log", "moment", "Statistics"];

    function StatisticsController($scope, $log, moment, Statistics) {
        $scope.stats = Statistics.get(
            {}, function(statsData) {
                $scope.retrieved = true;
            }
        );
    }
})();
