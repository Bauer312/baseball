package main

import (
	"fmt"
	"net/url"
	"time"
)

/*
GameURL is an interface to a URL-building mechanism for game data.
	The MLB website has a pretty nice URL structure; given some
	parameters a URL can be built so that resource can be retrieved.
*/
type GameURL interface {
	SetBaseURL(url string)
	GetURLsForDate(date time.Time) string
	GetURLsForDates(dates []time.Time) []*url.URL
}

/*
LocalURL is an implementation of the GameURL interface
*/
type LocalURL struct {
	baseURL string
}

/*
SetBaseURL specifies the baseURL to be used for constructing URLs
*/
func (lu *LocalURL) SetBaseURL(url string) {
	lu.baseURL = url
}

/*
GetURLsForDate returns a URL for a specific day
*/
func (lu *LocalURL) GetURLsForDate(date time.Time) string {
	year := date.Year()
	month := date.Month()
	day := date.Day()
	return fmt.Sprintf("%syear_%04d/month_%02d/day_%02d", lu.baseURL, year, month, day)
}

/*
GetURLsForDates returns a slice of all URLs that correspond to that date
*/
func (lu *LocalURL) GetURLsForDates(dates []time.Time) []*url.URL {
	returnUrls := make([]*url.URL, len(dates))
	for i, date := range dates {
		rawurl := lu.GetURLsForDate(date)
		newURL, err := url.Parse(rawurl)
		if err != nil {
			fmt.Printf("Error parsing raw URL: %s", err)
		}
		returnUrls[i] = newURL
		fmt.Println(rawurl)
	}
	return returnUrls
}
