package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/go-errors/errors"
	"github.com/heroku/rollrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	eventFilenameRE = regexp.MustCompile(
		`(\d{4})-(\d{2})-(\d{2})-(\d{1,2})`)
	// AppEnv is: production, staging, development
	AppEnv string
)

func init() {
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

// Digest contains all aggregate data for specific hour
// +gen * slice:"SortBy"
type Digest struct {
	Count int       `json:"count"`
	Date  time.Time `json:"date"`
}

// EventRecord is one transformed event
type EventRecord struct {
	Actor ActorRecord `json:"actor"`
}

// ActorRecord is often nested in EventRecord
type ActorRecord struct {
	Username string `json:"login"`
}

// Username implements set methods
// +gen set
type Username string

// DigestFile will return a valid Digest instance based on a file,
// using a cached digest if available
func DigestFile(eventFilePath string, users UsernameSet) (*Digest, error) {
	digestFilePath := fmt.Sprintf("%v.digest.json", eventFilePath)
	df, err := os.OpenFile(digestFilePath,
		os.O_EXCL|os.O_CREATE|os.O_RDWR,
		0664)
	if err != nil {
		if os.IsExist(err) {
			return readDigest(digestFilePath)
		}
		return nil, err
	}
	defer df.Close()

	return doDigestFile(eventFilePath, df, users)
}

func doDigestFile(eventFilePath string, digestFile *os.File,
	users UsernameSet) (*Digest, error) {
	f, err := os.Open(eventFilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader, err := gzip.NewReader(f)
	if err != nil {
		panic(err)
	}

	c, err := lineCounter(reader)
	if err != nil {
		panic(err)
	}
	if _, err := f.Seek(0, 0); err != nil {
		panic(err)
	}

	reader.Reset(f)

	err = usernameExtractor(reader, users)
	if err != nil {
		panic(err)
	}

	dateParts := eventFilenameRE.FindStringSubmatch(
		filepath.Base(eventFilePath))
	year, _ := strconv.Atoi(dateParts[1])
	month, _ := strconv.Atoi(dateParts[2])
	day, _ := strconv.Atoi(dateParts[3])
	hr, _ := strconv.Atoi(dateParts[4])
	fileDate := time.Date(
		year, time.Month(month), day, hr,
		0, 0, 0, time.UTC)
	digest := &Digest{
		Count: c,
		Date:  fileDate,
	}
	if err != nil {
		return nil, err
	}

	log.Debugf("computed %v: %v events\n", fileDate, c)
	err = json.NewEncoder(digestFile).Encode(digest)
	return digest, err
}

func readDigest(digestFilePath string) (*Digest, error) {
	f, err := os.Open(digestFilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	d := &Digest{}
	err = json.NewDecoder(f).Decode(d)
	return d, err
}

func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 1024*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return count, errors.New(err)
		}

		count += bytes.Count(buf[:c], lineSep)

		if err == io.EOF {
			break
		}
	}

	return count, nil
}

func usernameExtractor(r io.Reader, users UsernameSet) error {
	decoder := json.NewDecoder(r)
	for decoder.More() {
		event := EventRecord{}
		err := decoder.Decode(&event)
		if err != nil {
			return err
		}
		event.Actor.Username = strings.ToLower(event.Actor.Username)
		users.Add(Username(event.Actor.Username))
	}
	return nil
}

func makePath(basename string) string {
	return filepath.Join(
		os.Getenv("GHC_EVENTS_PATH"),
		basename)
}

func makeSummary(digests DigestSlice, newUsers UsernameSet) {
	digests = DigestSlice(digests).SortBy(func(x, y *Digest) bool {
		return x.Date.Unix() < y.Date.Unix()
	})

	digestSummary, err := os.Create(makePath("summary.json"))
	if err != nil {
		panic(err)
	}
	defer digestSummary.Close()

	err = json.NewEncoder(digestSummary).Encode(digests)
	if err != nil {
		panic(err)
	}

	usersSummary, err := os.OpenFile(
		makePath("users.txt"),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664)
	if err != nil {
		panic(err)
	}
	defer usersSummary.Close()

	log.Debugf("writing %v users\n", len(newUsers))
	for u := range newUsers {
		_, err = fmt.Fprintln(usersSummary, u)
		if err != nil {
			panic(err)
		}
	}
}

func readKnownUsers() UsernameSet {
	users := UsernameSet{}
	usersBuf, err := ioutil.ReadFile(makePath("users.txt"))
	if err == nil {
		userStrings := strings.Split(string(usersBuf), "\n")
		for _, u := range userStrings {
			users.Add(Username(u))
		}
	} else {
		log.WithError(err).Warn("warning: could not read users.txt")
	}
	return users
}

func main() {
	log.Debugf("event digest started")
	users := readKnownUsers()
	existingUsers := users.Clone()
	log.Debugf("found %v existing users", len(existingUsers))

	eventFiles, err := filepath.Glob(makePath("*.json.gz"))
	if err != nil {
		panic(err)
	}

	digests := make([]*Digest, 0, len(eventFiles))
	for _, f := range eventFiles {
		d, err := DigestFile(f, users)
		if err != nil {
			panic(err)
		}
		log.Debugf("now have %v users", len(users))
		digests = append(digests, d)
	}

	log.Debug("computing difference in users")
	newUsers := users.Difference(existingUsers)
	log.Debugf("done (found %v)\n", len(newUsers))
	makeSummary(digests, newUsers)
}
