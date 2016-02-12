(function () {
    angular
        .module("ghca")
        .directive("eventOcticon", EventOcticon);

    function EventOcticon() {
        var octiconMap = {
            "GollumEvent": "book",
            "IssuesEvent": "issue-opened",
            "PushEvent": "repo-push",
            "CommitCommentEvent": "comment",
            "ReleaseEvent": "tag",
            "PublicEvent": "megaphone",
            "MemberEvent": "person",
            "IssueCommentEvent": "comment-discussion"
        };

        var eventDescriptionMap = {
            "GollumEvent": "Wiki",
            "IssuesEvent": "Issue",
            "PushEvent": "Push",
            "CommitCommentEvent": "Commit Comment",
            "ReleaseEvent": "Release",
            "PublicEvent": "Repository made public",
            "MemberEvent": "Membership change",
            "IssueCommentEvent": "Issue comment"
        };

        return {
            restrict: "A",
            require: "^ngModel",
            scope: {
                ngModel: '='
            },
            template: '',
            link: function(scope, element, attrs) {
                element.addClass("octicon");
                element.addClass("octicon-" + octiconMap[scope.ngModel]);
                element.attr("data-toggle", "tooltip");
                element.attr("data-placement", "left");
                element.attr("title", eventDescriptionMap[scope.ngModel]);
                $(element).tooltip();
            }
        };
    }
})();
