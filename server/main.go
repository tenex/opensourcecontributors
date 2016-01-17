package main

import (
	"net/http"

	"encoding/json"

	"gopkg.in/mgo.v2"
)

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
