package main

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/go-errors/errors"
	"github.com/heroku/rollrus"
	"gopkg.in/natefinch/lumberjack.v2"
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
	logDest := os.Getenv("GHC_EVENT_DIGEST_LOG_PATH")
	if logDest == "" {
		logDest = "/var/log/ghc/event-digest.log"
	}
	if AppEnv == "development" {
		log.SetLevel(log.DebugLevel)
	}
	if AppEnv == "production" {
		log.SetOutput(&lumberjack.Logger{
			Filename: logDest,
			MaxSize:  100, // MB
		})
		rollrus.SetupLogging(os.Getenv("GHC_ROLLBAR_TOKEN"), AppEnv)
		// PUT THIS AFTER ROLLRUS!
		// https://github.com/heroku/rollrus/issues/4
		log.SetFormatter(&log.JSONFormatter{})
	}
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
			digest, readErr := readDigest(digestFilePath)
			if readErr == nil {
				return digest, nil
			}
			evt := log.WithField("eventFile", eventFilePath)
			evt.Warn("could not read existing digest")
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
		return nil, errors.New(err)
	}

	c, err := digestStream(reader, users)
	if err != nil {
		entry := log.WithError(err)
		entry = entry.WithField("eventFilePath", eventFilePath)
		entry.Error("could not digest stream")
		return nil, err
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

	log.Debugf("computed %v: %v events", fileDate, c)
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

// ErrNoNullBytes indicates that null bytes were expected but not found
var ErrNoNullBytes = errors.New("no null bytes found to skip")

// skipMysteriousNulls provides workaround for:
// https://github.com/igrigorik/githubarchive.org/issues/135
func skipMysteriousNulls(r *bufio.Reader) error {
	skippedBytes := false
	for {
		c, err := r.ReadByte()
		if err != nil {
			return err
		}
		if c != 0x00 {
			r.UnreadByte()
			if !skippedBytes {
				return ErrNoNullBytes
			}
			return nil
		}
		skippedBytes = true
	}
}

func digestStream(r io.Reader, users UsernameSet) (int, error) {
	records := 0
	bufReader := bufio.NewReader(r)
	decoder := json.NewDecoder(bufReader)
	for {
		event := EventRecord{}
		if err := decoder.Decode(&event); err == io.EOF {
			break
		} else if err != nil {
			// Could be because of mysterious nulls...
			nullErr := skipMysteriousNulls(bufReader)
			if nullErr != nil {
				entry := log.WithError(nullErr)
				entry.Warn("encounterd error while skipping nulls")
				return records, err
			}
			// json.Decoder stores its own buffer
			decoder = json.NewDecoder(bufReader)
			// try again
			continue
		}
		records++
		event.Actor.Username = strings.ToLower(event.Actor.Username)
		users.Add(Username(event.Actor.Username))
	}
	return records, nil
}

func makePath(basename string) string {
	return filepath.Join(
		os.Getenv("GHC_EVENTS_PATH"),
		basename)
}

func makeSummary(digests DigestSlice, newUsers UsernameSet) {
	digests = DigestSlice(digests).SortBy(func(x, y *Digest) bool {
		if x == nil || y == nil {
			ent := log.WithField("x", x).WithField("y", y)
			ent.Warn("encountered nil digest during sort")
			return false
		}
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

	log.Debugf("writing %v users", len(newUsers))
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
			// if a file can't be digested, then many problably
			// can't, so we should investigate
			log.WithError(err).Errorf(
				"could not digest events file")
			break
		}
		digests = append(digests, d)
	}

	log.Debug("computing difference in users")
	newUsers := users.Difference(existingUsers)
	log.Debugf("done (found %v new users)", len(newUsers))
	makeSummary(digests, newUsers)
}
