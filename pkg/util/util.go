/*
	Copyright 2017 Brian Bauer

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

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
DateToURLNoSideEffects will turn a date into a URL, given a rootURL
*/
func DateToURLNoSideEffects(date time.Time, root string) (*url.URL, error) {
	if len(root) == 0 {
		return nil, errors.New("The root URL has not been set")
	}
	year := date.Year()
	month := date.Month()
	day := date.Day()
	rawURL := fmt.Sprintf("%s/year_%04d/month_%02d/day_%02d/", root, year, month, day)

	realURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	return realURL, nil
}

/*
GameToURLsNoSideEffect will turn a game into a set of URLs, given a rootURL
*/
func GameToURLsNoSideEffect(game, root string) ([]*url.URL, error) {
	if len(root) == 0 {
		return nil, errors.New("The root URL has not been set")
	}
	year := game[4:8]
	month := game[9:11]
	day := game[12:14]
	rawURL := fmt.Sprintf("%s/year_%s/month_%s/day_%s/%s", root, year, month, day, game)
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
URLToFSPathNoSideEffects will turn a URL into a filesystem path.  If the URL doesn't specify an
	actual file name, append index.html to it because that is what the web server
	is going to do.  Also, get rid of some of the intermediate portions of the
	URL to prevent the filesystem path from being unnecessarily long.
*/
func URLToFSPathNoSideEffects(realURL *url.URL, root string) (string, error) {
	if len(root) == 0 {
		return "", errors.New("The root filesystem path has not been set")
	}

	rawString := realURL.Path
	switch {
	case strings.HasSuffix(rawString, ".html"):
	case strings.HasSuffix(rawString, ".xml"):
	case strings.HasSuffix(rawString, ".dat"):
	default:
		rawString = rawString + "index.html"
	}
	pathComponents := strings.Split(rawString, "/")
	pathComponents = pathComponents[4:]
	newPath := strings.Join(pathComponents, "/")
	newPath = filepath.Join(root, newPath)

	return newPath, nil
}

/*
DateToProcessedFileNoSideEffects will turn a date and a processed file name into a filesystem path
*/
func DateToProcessedFileNoSideEffects(date time.Time, rootURL, rootFS, processedFile string) (string, error) {
	dateURL, err := DateToURLNoSideEffects(date, rootURL)
	if err != nil {
		return "", err
	}
	newURL, err := url.Parse(dateURL.Path + processedFile)
	if err != nil {
		return "", err
	}
	dateFS, err := URLToFSPathNoSideEffects(newURL, rootFS)
	if err != nil {
		return "", err
	}
	return dateFS, nil
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
func SaveURLToPath(targetURL *url.URL, targetPath string, client *http.Client) error {
	// First, make sure the directory exists
	err := VerifyFSDirectory(targetPath)
	if err != nil {
		return err
	}
	fmt.Printf("Saving %s to %s\n", targetURL, targetPath)

	// Second, make the request
	res, err := client.Get(targetURL.String())
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
