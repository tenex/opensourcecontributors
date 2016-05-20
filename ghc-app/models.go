package main

import (
	"sort"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// PageSize set the size of the result set where applicable
const PageSize = 50

// UserContributionsFunc returns raw BSON records of a user's contributions
type UserContributionsFunc func(string, int) ([]bson.M, error)

// UserContributionsFactory returns a UserContributionsFunc that can be used
// to retrieve a user's contributions given their username and a zero-based skip
func UserContributionsFactory(c *mgo.Collection) UserContributionsFunc {
	return func(username string, skip int) ([]bson.M, error) {
		username = strings.ToLower(username)

		var events []bson.M
		query := c.Find(bson.M{"_user_lower": username})
		query = query.Sort("-created_at")
		query = query.Skip(skip).Limit(PageSize)
		err := query.All(&events)
		if err != nil {
			return nil, err
		}
		return events, nil
	}
}

//////////////////
// User Summary //
//////////////////

// UserSummary describes a brief summary of a user's contributions
type UserSummary struct {
	Username     string   `json:"username"`
	Repositories []string `json:"repos"`
	EventCount   int      `json:"eventCount"`
}

// UserSummaryFunc returns an instance of UserSummary
// given a username
type UserSummaryFunc func(string) (*UserSummary, error)

// UserSummaryFactory returns a UserSummaryFunc that can be used
// to retrieve a user's summary
func UserSummaryFactory(c *mgo.Collection) UserSummaryFunc {
	return func(username string) (*UserSummary, error) {
		username = strings.ToLower(username)
		query := c.Find(bson.M{"_user_lower": username})

		repoList := []string{}
		err := query.Distinct("repo", &repoList)
		if err != nil {
			return nil, err
		}
		sort.Strings(repoList)

		ct, err := query.Count()
		if err != nil {
			return nil, err
		}

		return &UserSummary{
			Username:     username,
			Repositories: repoList,
			EventCount:   ct,
		}, nil
	}
}

////////////////
// Statistics //
////////////////

// GHCStats describes statistics about the project's database
type GHCStats struct {
	EventCount     int       `json:"eventCount"`
	LatestEvent    time.Time `json:"latestEvent"`
	LatestEventAge int64     `json:"latestEventAge"`
}

// GHCStatsFunc returns the latest statistics
type GHCStatsFunc func() (*GHCStats, error)

// GHCStatsFactory returns a GHCStatsFunc which can be used
// to reteurn statistics about the project
func GHCStatsFactory(c *mgo.Collection) GHCStatsFunc {
	return func() (*GHCStats, error) {
		ct, err := c.Count()
		if err != nil {
			return nil, err
		}

		var latestEvt bson.M
		err = c.Find(nil).Sort("-created_at").One(&latestEvt)
		if err != nil {
			return nil, err
		}
		latestEvtTime, err := time.Parse(
			time.RFC3339,
			latestEvt["created_at"].(string))
		latestEvtAge := int64(
			time.Now().UTC().Sub(latestEvtTime).Seconds())

		return &GHCStats{
			EventCount:     ct,
			LatestEvent:    latestEvtTime,
			LatestEventAge: latestEvtAge,
		}, nil
	}

}

/////////////
// Summary //
/////////////

// GHCSummary respresents a collection of daily summaries
type GHCSummary struct {
	DailySummary []GHCDailySummary `json:"dailySummary"`
}

// GHCDailySummary represents one day's counts of data
type GHCDailySummary struct {
	Date  string `bson:"date" json:"date"`
	Count int64  `bson:"count" json:"count"`
}

// GHCSummaryFunc returns daily summaries of events
type GHCSummaryFunc func() (*GHCSummary, error)

// GHCSummaryFactory is great
func GHCSummaryFactory(c *mgo.Collection) GHCSummaryFunc {
	return func() (*GHCSummary, error) {
		summary := GHCSummary{}
		err := c.Find(nil).Sort("-date").All(&summary.DailySummary)
		if err != nil {
			return nil, err
		}
		return &summary, nil
	}
}
