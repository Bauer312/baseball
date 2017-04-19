package util

import (
	"errors"
	"fmt"
	"net/url"
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
	if strings.HasSuffix(fs, "/") == false {
		fs = fs + "/"
	}
	rootFS = fs
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
	newPath := rootFS + strings.Join(pathComponents, "/")

	return newPath, nil
}
