package main

import (
	"encoding/json"
	"net/http"

	"bytes"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
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
	controller.HandleFunc("/user/{username}/events", controller.UserEvents)
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

// UserEvents is a controller action for /user/{username}/events
func (c *GHCController) UserEvents(rw http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	contributionBSON, err := c.userContributions(username)
	if err != nil {
		panic(err) // TODO: Fix
	}
	err = serveJSON(rw, contributionBSON)
	if err != nil {
		panic(err) // TODO: Fix
	}
}

// UserSummary is a controller action for /user/{username}
func (c *GHCController) UserSummary(rw http.ResponseWriter, r *http.Request) {
	// TODO
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
