package main

import (
	"strings"

	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// UserContributionsFunc returns raw BSON records of a user's contributions
type UserContributionsFunc func(string) ([]bson.M, error)

// UserContributionsFactory returns a UserContributionsFunc that can be used
// to retrieve a user's contributions given their username
func UserContributionsFactory(c *mgo.Collection) UserContributionsFunc {
	return func(username string) ([]bson.M, error) {
		username = strings.ToLower(username)
		var maps []bson.M
		err := c.Find(
			bson.M{"_user_lower": username},
		).All(&maps)
		if err != nil {
			return nil, err
		}
		return maps, nil
	}
}

//////////////////
// User Summary //
//////////////////

// UserSummary describes a brief summary of a user's contributions
type UserSummary struct {
	Username          string
	Repositories      []string
	ContributionCount int
}

// UserSummaryFunc returns an instance of UserSummary
// given a username
type UserSummaryFunc func(string) (*UserSummary, error)

// UserSummaryFactory returns a UserSummaryFunc that can be used
// to retrieve a user's summary
func UserSummaryFactory(c *mgo.Collection) UserSummaryFunc {
	return func(username string) (*UserSummary, error) {
		username = strings.ToLower(username)
		// summary = {
		// 	"username": username,
		// 	"eventCount": event_count,
		// 	"repos": repos,
		// }
		return &UserSummary{
			Username: username,
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
		// TODO: Implement
		return &GHCStats{}, nil
	}

}
