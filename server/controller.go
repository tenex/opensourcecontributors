package main

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
)

// GHCController is the master controller for this application
type GHCController struct {
	// Contributions : MongoDB contributions.contributions
	Contributions *mgo.Collection
	*mux.Router

	userContributions UserContributionsFunc
}

func NewGHCController(contributions *mgo.Collection) *GHCController {
	controller := &GHCController{
		Contributions:     contributions,
		userContributions: UserContributionsFactory(contributions),
	}
	router := mux.NewRouter()
	router.Handle("/", controller)
	router.HandleFunc("/user/{username}", controller.UserSummary)
	router.HandleFunc("/user/{username}/events", controller.UserEvents)
	router.HandleFunc("/stats", controller.Stats)
	return controller
}

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

func (c *GHCController) UserSummary(rw http.ResponseWriter, r *http.Request) {
	// TODO
}

func (c *GHCController) Stats(rw http.ResponseWriter, r *http.Request) {
	n, err := c.Contributions.Count()
	if err != nil {
		panic(err)
	}
	rw.Write([]byte(strconv.Itoa(n)))
}
