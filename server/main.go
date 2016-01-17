package main

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"gopkg.in/mgo.v2"
)

func main() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	collection := session.DB("contributions").C("contributions")

	controller := NewGHCController(collection)
	loggedHandler := handlers.CombinedLoggingHandler(os.Stdout, controller)
	recoveryHandler := handlers.RecoveryHandler()(loggedHandler)
	http.ListenAndServe(":5000", recoveryHandler)
}
