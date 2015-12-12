#!/bin/bash
HOST="http://githubcontributions.io/stats"
MAX_AGE="7300" # seconds
SCAPEGOAT="liambowen@gmail.com"

function notify_failure()
{
    if [ -t 1 ]
    then
        echo "old records detected. notifying ${SCAPEGOAT}"
    fi
    age=$1
    msg="Error: recent events are not being entered into the database.\n"
    msg="${msg} Maximum age of newest record is ${MAX_AGE} seconds\n"
    msg="${msg} Currently the newest record is ${age} seconds old.\n"
    printf "${msg}" | mail -s "ERROR: GitHub Contributions" "${SCAPEGOAT}"
}

age="$(curl "${HOST}" 2>/dev/null | jq '.latestEventAge')"
(($age > $MAX_AGE)) && notify_failure "${age}"
