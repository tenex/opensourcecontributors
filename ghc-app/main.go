package main

//go:generate go-bindata -prefix static/ static/...

import (
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/heroku/rollrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	// AppEnv is: production, staging, development
	AppEnv string
)

func init() {
	// This is used to generate Request IDs
	rand.Seed(time.Now().UnixNano())

	AppEnv = os.Getenv("GHC_ENV")
	if AppEnv == "" {
		AppEnv = "development"
	}
	logDest := os.Getenv("GHC_APP_LOG_PATH")
	if logDest == "" {
		logDest = "/var/log/ghc/ghc.log"
	}
	log.SetOutput(&lumberjack.Logger{
		Filename: logDest,
		MaxSize:  100, // MB
	})
	if AppEnv == "production" {
		rollrus.SetupLogging(os.Getenv("GHC_ROLLBAR_TOKEN"), AppEnv)
	}
	// PUT THIS AFTER ROLLRUS!
	// https://github.com/heroku/rollrus/issues/4
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

// Stops panics with no panic-worthy cause
// Stops:
//   - EPIPE, which occurs when a client stops loading a page
func xanax(v interface{}) error {
	if v == nil {
		return nil
	}
	var err error
	switch cause := v.(type) {
	case error:
		err = cause
		if err == syscall.EPIPE {
			err = nil
		}
	default:
		err = fmt.Errorf("%v", v)
	}
	return err
}

func recoverHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := xanax(recover()); err != nil {
				rawStr := fmt.Sprintf("%#v", err)
				log.WithField("raw", rawStr).Error(err.Error())
				http.Error(rw, rawStr, 500)
			}
		}()
		h.ServeHTTP(rw, r)
	})
}

func remoteAddrHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		remoteAddr := r.Header.Get("X-Forwarded-For")
		if remoteAddr != "" {
			r.RemoteAddr = remoteAddr
		}
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
	session.SetSafe(nil)                       // we never write
	session.SetSocketTimeout(15 * time.Minute) // cheaper than SSDs
	session.SetMode(mgo.Monotonic, true)
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	bindEndpoint := net.JoinHostPort("", port)

	handler := mainHandler(session)
	handler = recoverHandler(handler)
	handler = logHandler(handler)
	handler = remoteAddrHandler(handler)

	http.ListenAndServe(bindEndpoint, handler)
}
