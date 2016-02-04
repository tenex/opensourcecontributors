(function() {
    angular
        .module('ghca')
        .controller("UserController",
                    UserController);

    UserController.$inject = [
        "$log", "User", "username"
    ];

    function UserController($log, User, username) {
        var vm = this;

        vm.username = username;
        vm.userUrl = "";

        vm.user = User.get({username: username}, function(user) {
            vm.userUrl = "https://github.com/" + user.username;
        });


        $log.debug("UserController: ", username);
    }
})();
