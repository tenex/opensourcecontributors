(function() {
    angular
        .module('ghca')
        .controller("UserRepositoriesController",
                    UserRepositoriesController);

    UserRepositoriesController.$inject = [
        "$log", "moment", "User", "Event"
    ];

    function UserRepositoriesController(user) {

    }
})();
