package main

import (
	"bytes"
	"encoding/json"
	"net/http"
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
	controller := &GHCController{
		userContributions: UserContributionsFactory(contributions),
		userSummary:       UserSummaryFactory(contributions),
		ghcStats:          GHCStatsFactory(contributions),
		Router:            mux.NewRouter(),
	}
	controller.HandleFunc("/", controller.serveRoot)

	controller.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/",
			http.HandlerFunc(controller.serveStatic)))
	controller.HandleFunc("/user/{username}", controller.UserSummary)
	controller.HandleFunc("/user/{username}/events",
		controller.UserEvents)
	controller.HandleFunc("/user/{username}/events/{page:[0-9]+}",
		controller.UserEvents)

	controller.HandleFunc("/stats", controller.Stats)
	return controller
}

func (c *GHCController) serveRoot(rw http.ResponseWriter, _ *http.Request) {
	content, _ := Asset("index.html")
	rw.Write(content)
}

func (c *GHCController) serveStatic(rw http.ResponseWriter, r *http.Request) {
	fi, err := AssetInfo(r.URL.Path)
	if err != nil {
		panic(err)
	}
	assetReader := bytes.NewReader(MustAsset(r.URL.Path))
	http.ServeContent(rw, r,
		fi.Name(),
		fi.ModTime(),
		assetReader)
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
