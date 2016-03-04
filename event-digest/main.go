package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

var (
	eventFilenameRE = regexp.MustCompile(
		`(\d{4})-(\d{2})-(\d{2})-(\d{1,2})`)
)

// Digest contains all aggregate data for specific hour
// +gen * slice:"SortBy"
type Digest struct {
	Count int        `json:"count"`
	Date  time.Time  `json:"date"`
	Users []Username `json:"users"`
}

// EventRecord is one transformed event
type EventRecord struct {
	Username Username `json:"_user_lower"`
}

// Username implements set methods
// +gen set
type Username string

// DigestFile will return a valid Digest instance based on a file,
// using a cached digest if available
func DigestFile(eventFilePath string) (*Digest, error) {
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

	return doDigestFile(eventFilePath, df)
}

func doDigestFile(eventFilePath string, digestFile *os.File) (*Digest, error) {
	f, err := os.Open(eventFilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	reader, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}

	c, err := lineCounter(reader)
	f.Seek(0, os.SEEK_SET)
	reader.Reset(f)

	usernames, err := usernameExtractor(reader)
	if err != nil {
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
		Users: usernames.ToSlice(),
	}
	if err != nil {
		return nil, err
	}

	fmt.Printf("computed %v: %v\n", fileDate, c)
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
			return count, err
		}

		count += bytes.Count(buf[:c], lineSep)

		if err == io.EOF {
			break
		}
	}

	return count, nil
}

func usernameExtractor(r io.Reader) (UsernameSet, error) {
	return nil, nil
}

func makePath(basename string) string {
	return filepath.Join(
		os.Getenv("GHC_EVENTS_PATH"),
		basename)
}

func makeSummary(digests DigestSlice) {
	digests = DigestSlice(digests).SortBy(func(x, y *Digest) bool {
		return x.Date.Unix() < y.Date.Unix()
	})

	digestSummary, err := os.Create(
		makePath("summary.json"))
	if err != nil {
		panic(err)
	}

	err = json.NewEncoder(digestSummary).Encode(digests)
	if err != nil {
		panic(err)
	}
}

func main() {
	eventFiles, err := filepath.Glob(
		makePath("*.json.gz"))
	if err != nil {
		panic(err)
	}

	digests := make([]*Digest, 0, len(eventFiles))

	for _, f := range eventFiles {
		d, err := DigestFile(f)
		if err != nil {
			panic(err)
		}
		digests = append(digests, d)
	}

	makeSummary(digests)
}
