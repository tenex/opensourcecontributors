package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// GHCController is the master controller for this application
type GHCController struct {
	*mux.Router

	userContributions UserContributionsFunc
	userSummary       UserSummaryFunc
	ghcStats          GHCStatsFunc
}

// NewGHCController is the constructor for GHCController
func NewGHCController(contributions *mgo.Collection) *GHCController {
	c := &GHCController{
		userContributions: UserContributionsFactory(contributions),
		userSummary:       UserSummaryFactory(contributions),
		ghcStats:          GHCStatsFactory(contributions),
		Router:            mux.NewRouter(),
	}
	c.HandleFunc("/user/{username}", c.UserSummary)
	c.HandleFunc("/user/{username}/events", c.UserEvents)
	c.HandleFunc("/user/{username}/events/{page:[0-9]+}", c.UserEvents)
	c.HandleFunc("/stats", c.Stats)
	c.HandleFunc("/error", c.Error)
	c.HandleFunc("/aggregates", c.Aggregates)
	return c
}

func (c *GHCController) Error(_ http.ResponseWriter, _ *http.Request) {
	panic(errors.New("error successful"))
}

// Aggregates serves /aggregates
func (c *GHCController) Aggregates(rw http.ResponseWriter, r *http.Request) {
	summaryPath := filepath.Join(
		os.Getenv("GHC_EVENTS_PATH"), "summary.json")
	f, err := os.Open(summaryPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	rw.WriteHeader(http.StatusOK)
	_, err = io.Copy(rw, f)
}

// UserEventsPage includes <= PageSize number of events and metadata about
// all of the events corresponding to the user
type UserEventsPage struct {
	Events      []bson.M `json:"events"`
	Start       int      `json:"start"`
	End         int      `json:"end"`
	CurrentPage int      `json:"currentPage"`
	PageCount   int      `json:"size"`
}

// UserEvents is a controller action for:
// /user/{username}/events
// /user/{username}/events/[page]
func (c *GHCController) UserEvents(rw http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	username := v["username"]
	page, err := strconv.Atoi(v["page"])
	if err != nil {
		page = 1
	}
	skip := (page - 1) * PageSize
	contributionBSON, err := c.userContributions(username, skip)
	if err != nil {
		panic(err) // TODO: Fix
	}
	eventsPage := UserEventsPage{
		Events:      contributionBSON,
		Start:       skip,
		End:         len(contributionBSON) + skip,
		CurrentPage: page,
		PageCount:   len(contributionBSON),
	}
	err = serveJSON(rw, eventsPage)
	if err != nil {
		panic(err) // TODO: Fix
	}
}

// UserSummary is a controller action for /user/{username}
func (c *GHCController) UserSummary(rw http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	summary, err := c.userSummary(username)
	if err != nil {
		panic(err) // TODO: Fix
	}
	err = serveJSON(rw, summary)
	if err != nil {
		panic(err) // TODO: Fix
	}
}

// Stats is a controller action for /stats
func (c *GHCController) Stats(rw http.ResponseWriter, r *http.Request) {
	stats, err := c.ghcStats()
	if err != nil {
		panic(err)
	}
	serveJSON(rw, stats)
}

func serveJSON(rw http.ResponseWriter, obj interface{}) error {
	return json.NewEncoder(rw).Encode(obj)
}
