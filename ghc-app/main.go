package main

//go:generate go-bindata -prefix static/ static/...

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	logDest := os.Getenv("GHC_APP_LOG_PATH")
	if logDest == "" {
		logDest = "/var/log/ghc/ghc.log"
	}
	log.SetOutput(&lumberjack.Logger{
		Filename: logDest,
		MaxSize:  100, // MB
	})
	log.SetFormatter(&log.JSONFormatter{})
}

func makeRequestID() string {
	return fmt.Sprintf("%08X", rand.Uint32())
}

func logHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rID := makeRequestID()
		log.WithFields(log.Fields{
			"requestID":  rID,
			"referer":    r.Referer(),
			"remoteAddr": r.RemoteAddr,
			"url":        r.URL.String(),
			"userAgent":  r.UserAgent(),
			"method":     r.Method,
		}).Info("request")
		w.Header().Add("X-GHC-Request-ID", rID)

		startTime := time.Now()
		h.ServeHTTP(w, r)
		elapsedTime := time.Now().Sub(startTime).Seconds()

		log.WithFields(log.Fields{
			"requestID": rID,
			"elapsed":   elapsedTime,
		}).Info("response")
	})
}

func recoverHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.WithField("error", err).Error("panic")
				http.Error(rw, fmt.Sprintf("%#v", err), 500)
			}
		}()
		h.ServeHTTP(rw, r)
	})
}

func mainHandler(globalSession *mgo.Session) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		session := globalSession.Copy()
		defer session.Close()
		collection := session.DB("contributions").C("contributions")
		NewGHCController(collection).ServeHTTP(rw, r)
	})
}

func main() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetSafe(nil) // we never write
	session.SetMode(mgo.Monotonic, true)

	handler := logHandler(recoverHandler(mainHandler(session)))
	http.ListenAndServe(":5000", handler)
}
