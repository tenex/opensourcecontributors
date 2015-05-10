# GitHub Contributions

This is a utility to find a list of all contributions a user has made to any public repository on GitHub from 2011-01-01 through yesterday.

The data from 2015-01-01 - present is found on [GitHub Archive](https://www.githubarchive.org). The data from before this uses a different schema and was obtained from Google's BigQuery (see below)

### Processing data archives

The tool `archive-processor` will transform either the Timeline or Event API archives into JSON files which can be imported directly into your database of choice. To generate the output files (which will be gzipped if given gzipped input), use:

```sh
 ./archive-processor --output-path /github-archive/processed --timeline-path /github-archive/2011-2014/ --events-path /github-archive/2015/
```

Logs will be dumped to your PWD and STDERR. It would be nice to refine that a bit. If you don't specify an `--output-path`, the processed JSON will be dumped to your STDOUT.

#### Time / Money

On my [vultr.com](http://www.vultr.com/?ref=6831514) VPS, with the cheapest plan at US $5 per month, it takes about 20 hours to process all this data with one script working on the timeline archives and a second one working on the event archives. I'm using a non-SSD disk because of the massive price difference. With the SSD, it would take about 8.3 hours for the timeline archive and about 4 hours for all the events archive data for 5 months.

## BigQuery Data Sets

For the data from 2011-2014 (actually, 2008-08-25 01:07:06 to 2014-12-31 23:59:59), the GitHub Archive project recorded data from the (now deprecated) Timeline API. This is in a different format and has many more quirks than the new [GitHub Events API](https://developer.github.com/v3/activity/events/). To obtain this data, the following BigTable query was used (which took only 47.5s to run):

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

#### PushEvent with no repository name (Timeline API)

Example:

```sql
SELECT *
FROM [githubarchive:year.2014]
WHERE payload_head='8824ed4d86f587a2a556248d9abfac790a1cbd3f'
LIMIT 1
```

It seems like sometimes, the only way to get the real repository name (`owner/project`) is to parse it from the URL.

#### PushEvent with no way of figuring out the repository (Timeline API)

Example:

```sql
SELECT *
FROM [githubarchive:year.2011]
WHERE payload_head='32b2177f05be005df3542c14d9a9985be2b553f7'
LIMIT 5
```

`repository_url` is `https://github.com//` and `repository_name` is `/` for each of these. They actually push to:
https://github.com/Jiyambi/WoW-Pro-Guides but I only know that by reading the commit messages.
