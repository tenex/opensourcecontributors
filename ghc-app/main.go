package main

//go:generate go-bindata -prefix static/ static/...

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/thoas/stats"
	"gopkg.in/mgo.v2"
	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	logDest := os.Getenv("GHC_APP_LOG_PATH")
	if logDest == "" {
		logDest = "/var/log/ghc/ghc.log"
	}
	fmt.Printf("logging to %s\n", logDest)
	log.SetOutput(&lumberjack.Logger{
		Filename: logDest,
		MaxSize:  100, // MB
	})
}

func printStatsLoop(s *stats.Stats) {
	for {
		fmt.Printf("%#v\n", s.Data())
		time.Sleep(30 * time.Second)
	}

}

func makeRequestID() string {
	return fmt.Sprintf("%08X", rand.Uint32())
}

func logHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rID := makeRequestID()
		log.WithFields(log.Fields{
			"requestID":  rID,
			"referer":    r.Referer(),
			"remoteAddr": r.RemoteAddr,
			"url":        r.URL.String(),
			"userAgent":  r.UserAgent(),
			"method":     r.Method,
		}).Info("request")
		startTime := time.Now()
		fn(w, r)
		elapsedTime := time.Now().Sub(startTime).Seconds()
		log.WithFields(log.Fields{
			"requestID": rID,
			"elapsed":   elapsedTime,
		}).Info("response")
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
