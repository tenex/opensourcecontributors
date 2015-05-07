# GitHub Contributions

This is a utility to find a list of all contributions a user has made to any public repository on GitHub from 2011-01-01 through yesterday.

The data from 2015-01-01 is found on [GitHub Archive](https://www.githubarchive.org). The data from before this uses a different schema and was obtained from Google's BigQuery (see below)

Place all these files from 2015-01-01 until today in a directory pointed to by the environment variable `ARCHIVE_PATH`, or in `~/github-archive`.

## BigQuery Data Sets

For the data from 2011-2014, this BigTable query was used:

```sql
SELECT
  -- common fields
  created_at, actor, repository_owner, repository_name, repository_organization, type, url,
  -- specific to type
  payload_action,            -- MemberEvent (added), IssuesEvent (open/closed)
  payload_member_login,      -- MemberEvent
  payload_commit_msg,        -- PushEvent
  payload_commit_email,      -- PushEvent
  payload_commit_id,         -- PushEvent
  payload_head,              -- PushEvent
  payload_ref,               -- PushEvent
  payload_comment_url,       -- CommitCommentEvent
  payload_comment_commit_id, -- CommitCommentEvent
  payload_comment_position,  -- CommitCommentEvent
  payload_comment_path,      -- CommitCommentEvent
  payload_comment_body,      -- CommitCommentEvent
  repository_size,           -- PublicEvent
  payload_issue_id,          -- IssueCommentEvent
  payload_comment_id         -- IssueCommentEvent
FROM (
  TABLE_QUERY(githubarchive:year,'true') -- All the years!
)
WHERE type IN (
  "GollumEvent",
  "IssuesEvent",
  "PushEvent",
  "CommitCommentEvent",
  "ReleaseEvent",
  "PublicEvent",
  "MemberEvent",
  "IssueCommentEvent"
)

```

If you actually want to use this data, there's no need to run that query; just ask me for the CSVs. When gzipped, they are about 19GB.
