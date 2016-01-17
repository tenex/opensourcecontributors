package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/thoas/stats"
	"gopkg.in/mgo.v2"
	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {
	log.SetOutput(&lumberjack.Logger{
		Filename: "/var/log/ghc/ghc.log",
		MaxSize:  100, // MB
	})
}

func printStatsLoop(s *stats.Stats) {
	for {
		fmt.Printf("%v\n", s.Data())
		time.Sleep(30 * time.Second)
	}

}

func logHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Before")
		fn(w, r)
		log.Println("After")
	}
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

	s := stats.New()
	go printStatsLoop(s)
	http.ListenAndServe(":5000",
		logHandler(
			s.Handler(controller).ServeHTTP))
}
