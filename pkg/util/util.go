package util

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

/*
Given a date, turn it into a URL.  Given a URL, turn it into a
	local filesytem path.
*/

var (
	rootURL string
	rootFS  string
)

/*
SetRoot is used to set the roots that will be used for constructing
	actual URLs and Filesystem paths
*/
func SetRoot(url, fs string) {
	rootURL = url
	rootFS = filepath.Clean(fs)
}

/*
DateToURL will turn a date into a URL
*/
func DateToURL(date time.Time) (*url.URL, error) {
	if len(rootURL) == 0 {
		return nil, errors.New("The root URL has not been set")
	}
	year := date.Year()
	month := date.Month()
	day := date.Day()
	rawURL := fmt.Sprintf("%s/year_%04d/month_%02d/day_%02d/", rootURL, year, month, day)

	realURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	return realURL, nil
}

/*
GameToURLs will turn a game into a set of URLs
*/
func GameToURLs(game string) ([]*url.URL, error) {
	if len(rootURL) == 0 {
		return nil, errors.New("The root URL has not been set")
	}
	year := game[4:8]
	month := game[9:11]
	day := game[12:14]
	rawURL := fmt.Sprintf("%s/year_%s/month_%s/day_%s/%s", rootURL, year, month, day, game)
	gameURLs := make([]*url.URL, 4)
	newURL, err := url.Parse(rawURL + "game.xml")
	if err != nil {
		return nil, err
	}
	gameURLs[0] = newURL
	newURL, err = url.Parse(rawURL + "game_events.xml")
	if err != nil {
		return nil, err
	}
	gameURLs[1] = newURL
	newURL, err = url.Parse(rawURL + "inning/inning_all.xml")
	if err != nil {
		return nil, err
	}
	gameURLs[2] = newURL
	newURL, err = url.Parse(rawURL + "inning/inning_hit.xml")
	if err != nil {
		return nil, err
	}
	gameURLs[3] = newURL
	return gameURLs, nil
}

/*
URLToFSPath will turn a URL into a filesystem path.  If the URL doesn't specify an
	actual file name, append index.html to it because that is what the web server
	is going to do.  Also, get rid of some of the intermediate portions of the
	URL to prevent the filesystem path from being unnecessarily long.
*/
func URLToFSPath(realURL *url.URL) (string, error) {
	if len(rootFS) == 0 {
		return "", errors.New("The root filesystem path has not been set")
	}

	rawString := realURL.Path
	switch {
	case strings.HasSuffix(rawString, ".html"):
	case strings.HasSuffix(rawString, ".xml"):
	default:
		rawString = rawString + "index.html"
	}
	pathComponents := strings.Split(rawString, "/")
	pathComponents = pathComponents[4:]
	newPath := strings.Join(pathComponents, "/")
	newPath = filepath.Join(rootFS, newPath)

	return newPath, nil
}

/*
VerifyFSDirectory verifies that a directory exists on the system, and creates it if it doesn't exist.
*/
func VerifyFSDirectory(fsPath string) error {
	fsDir := filepath.Dir(fsPath)
	return os.MkdirAll(fsDir, 0777)
}

/*
SaveURLToPath downloads a URL to a specific path on the filesystem
*/
func SaveURLToPath(targetURL *url.URL, targetPath string) error {
	// First, make sure the directory exists
	err := VerifyFSDirectory(targetPath)
	if err != nil {
		return err
	}
	fmt.Printf("Saving %s to %s\n", targetURL, targetPath)

	// Second, make the request
	res, err := http.Get(targetURL.String())
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Third, create the target file
	filePtr, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer filePtr.Close()

	// Fourth, copy the data into the file
	_, err = io.Copy(filePtr, res.Body)
	if err != nil {
		return err
	}

	return nil
}
