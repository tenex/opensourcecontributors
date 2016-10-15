# OpenSourceContributo.rs

[OpenSourceContributo.rs](http://opensourcecontributo.rs)

**Note about name change: This project was formerly known as githubcontributions.io. GitHub requested that the name of the project be changed in order to avoid confusion about who owns and maintains this project.**

This is a utility to find a list of all contributions a user has made to any public repository on GitHub from 2011-01-01 through yesterday.

The data from 2015-01-01 - present is found on [GitHub Archive](https://www.githubarchive.org). The data from before this uses a different schema and was obtained from Google's BigQuery (see below)

As of 2015-08-28, it tracks a total of
```sh
% cd /github-archive/processed
% gzip -l *.json.gz | awk 'END{print $2}' | numfmt --to=iec-i --suffix=B --format="%3f"
93GiB
% zcat *.json.gz | wc -l
253027947
```
events.

`db.contributions.stats()`:

```json
{
  "ns" : "contributions.contributions",
  "count" : 284048099,
  "size" : 113714359272,
  "avgObjSize" : 400,
  "storageSize" : 47820357632,
  "capped" : false,
  "nindexes" : 4,
  "totalIndexSize" : 8810385408,
  "indexSizes" : {
    "_id_" : 2804744192,
    "_user_lower_1" : 2275647488,
    "_event_id_1" : 1029251072,
    "created_at_1" : 2700742656
  },
  "ok" : 1
}
```
(WiredTiger stats omitted)

### Processing data archives

Processing the data archives involves 3 steps:

1. Download the raw events files from [GitHub Archive](https://www.githubarchive.org) into the events directory
2. Transform the events files by filtering non-contribution events (e.g., starring a repository) and adding necessary indexable keys (e.g., lowercased username)
3. Load the transformed data into MongoDB

The `archive-processor` tool in the `util` directory handles all of this.

The transformed data from step 2 is compressed and saved just in case we need to re-load the entire database (these files are much smaller than the raw data).

All of this can be done automatically by setting the correct environment variables, then running `archive-processor process`, or it can be invoked differently to separate the steps or change the working directories. Run `archive-processor --help` for details.

| Environment Variable | Meaning
|----------------------|----------------------------------------------------------|
| GHC_EVENTS_PATH      | Contains data from 2015-01-01 to present (.json.gz)      |
| GHC_TIMELINE_PATH    | Contains data before 2015-01-01 (.csv.gz)                |
| GHC_TRANSFORMED_PATH | Contains output of "transform" operation (.json.gz)      |
| GHC_LOADED_PATH      | Links to files in GHC_TRANSFORMED_PATH when loaded to DB |
| GHC_LOG_PATH         | Each invocation of `archive-processor` logs to here      |


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

## Credits

Created by [@hut8](http://github.com/hut8) and maintained by [Tenex Developers](https://tenex.tech) ([@tenex](http://github.com/tenex)).
