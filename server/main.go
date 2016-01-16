package main

import (
	"net/http"

	"strings"

	"encoding/json"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type UserContributionsFunc func(string) ([]bson.M, error)

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

type UserSummary struct {
	Username          string
	Repositories      []string
	ContributionCount int
}

type UserSummaryFunc func(string) (*UserSummary, error)

func UserSummaryFactory(c *mgo.Collection) UserSummaryFunc {
	return func(username string) (*UserSummary, error) {
		username = strings.ToLower(username)

		return &UserSummary{
			Username: username,
		}, nil
	}
}

func serveJSON(rw http.ResponseWriter, obj interface{}) error {
	return json.NewEncoder(rw).Encode(obj)
}

func main() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	collection := session.DB("contributions").C("contributions")
	controller := NewGHCController(collection)
	http.ListenAndServe(":5000", controller)
}
