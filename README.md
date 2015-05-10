# GitHub Contributions

This is a utility to find a list of all contributions a user has made to any public repository on GitHub from 2011-01-01 through yesterday.

The data from 2015-01-01 is found on [GitHub Archive](https://www.githubarchive.org). The data from before this uses a different schema and was obtained from Google's BigQuery (see below)

Place all these files from 2015-01-01 until today in a directory pointed to by the environment variable `ARCHIVE_PATH`, or in `~/github-archive`.

## BigQuery Data Sets

For the data from 2011-2014 (actually, 2008-08-25 01:07:06 to 2014-12-31 23:59:59), this BigTable query was used (which took me 47.5s to run):

```sql
SELECT
  -- common fields
  created_at, actor, repository_owner, repository_name, repository_organization, type, url,
  -- specific to type
  payload_page_html_url,     -- GollumEvent
  payload_page_summary,      -- GollumEvent
  payload_page_page_name,    -- GollumEvent
  payload_page_action,       -- GollumEvent
  payload_page_title,        -- GollumEvent
  payload_page_sha,          -- GollumEvent
  payload_number,            -- IssuesEvent
  payload_action,            -- MemberEvent, IssuesEvent, ReleaseEvent, IssueCommentEvent
  payload_member_login,      -- MemberEvent
  payload_commit_msg,        -- PushEvent
  payload_commit_email,      -- PushEvent
  payload_commit_id,         -- PushEvent
  payload_head,              -- PushEvent
  payload_ref,               -- PushEvent
  payload_comment_commit_id, -- CommitCommentEvent
  payload_comment_path,      -- CommitCommentEvent
  payload_comment_body,      -- CommitCommentEvent
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

### Erroneous data

There is lots of data in the archives that just doesn't make sense. Where I can, I've worked around it, for example by parsing needed data out of the event's URL. Here are some issues:

#### BigQuery exports CSV nulls weird?

Example:

```sql
SELECT *
FROM [githubarchive:year.2014]
LIMIT 1000
```

you will note that in the results pane of Google's BigQuery page, there is the string "null" where it really means a real null value. That makes its way into the exported CSV. So you should export the table the real way, or you will have the string "null" for almost every value.

#### PushEvent with no way of figuring out the repository (Timeline)

Example:

```sql
SELECT *
FROM [githubarchive:year.2011]
WHERE payload_head='32b2177f05be005df3542c14d9a9985be2b553f7'
LIMIT 5
```

`repository_url` is `https://github.com//` and `repository_name` is `/` for each of these. They actually push to:
https://github.com/Jiyambi/WoW-Pro-Guides but I only know that by reading the commit messages.
