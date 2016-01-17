package main

import (
	"github.com/codegangsta/negroni"
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
	n := negroni.Classic()
	n.UseHandler(controller)
	n.Run(":5000")
}
