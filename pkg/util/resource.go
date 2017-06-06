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
	"net/url"
	"time"
)

/*
TransferDefinition is a struct that contains the remote URL and the local filesystem locations for some resource
*/
type TransferDefinition struct {
	Source *url.URL
	Target string
}

type resourceConstructor interface {
	Roots(string, string) error
	Date(time.Time) ([]TransferDefinition, error)
	Game(string) ([]TransferDefinition, error)
}

/*
Resource implements the resourceConstructor interface, which enables us to turn dates and games into URLs and paths
*/
type Resource struct {
	rootURL string
	rootFS  string
}

/*
Roots is a method that receives the root URL and filesystem path and uses them when constructing URLs and paths
*/
func (r *Resource) Roots(url, fs string) error {
	if len(url) == 0 {
		return errors.New("url is not set")
	}
	r.rootURL = url

	if len(fs) == 0 {
		return errors.New("fs is not set")
	}
	r.rootFS = fs

	return nil
}

/*
Date is a method that turns the provided date into a slice containing a single transfer definition containing
    the remote URL and the local filesystem path for that date
*/
func (r *Resource) Date(date time.Time) ([]TransferDefinition, error) {
	if len(r.rootURL) == 0 {
		return nil, errors.New("The root URL has not been set")
	}

	dateURL, err := DateToURLNoSideEffects(date, r.rootURL)

	if err != nil {
		return nil, err
	}

	gameTransfers := make([]TransferDefinition, 1)
	gameTransfers[0].Source = dateURL
	gameTransfers[0].Target, err = URLToFSPathNoSideEffects(dateURL, r.rootFS)

	if err != nil {
		return nil, err
	}

	return gameTransfers, nil
}

/*
Game is a method that turns a game ID into a slice containing the transfer definitions for
	each desired file for that particular game
*/
func (r *Resource) Game(game string) ([]TransferDefinition, error) {
	if len(r.rootURL) == 0 {
		return nil, errors.New("The root URL has not been set")
	}

	gameURLs, err := GameToURLsNoSideEffect(game, r.rootURL)

	if err != nil {
		return nil, err
	}

	gameTransfers := make([]TransferDefinition, len(gameURLs))
	for i, gameURL := range gameURLs {
		gameTransfers[i].Source = gameURL
		gameTransfers[i].Target, err = URLToFSPathNoSideEffects(gameURL, r.rootFS)

		if err != nil {
			return nil, err
		}
	}

	return gameTransfers, nil
}
