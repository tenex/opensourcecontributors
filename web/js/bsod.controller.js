(function() {
    angular
        .module('ghca')
        .controller('BsodInstanceCtrl', BsodInstanceController);

    BsodInstanceController.$inject = [
        '$scope', '$uibModalInstance', 'errorDescription'
    ];

    function BsodInstanceController($scope, $uibModalInstance, errorDescription) {
        $scope.errorDescription = errorDescription;
        $scope.ok = function() {
            $uibModalInstance.close();
        };
    }
})();
