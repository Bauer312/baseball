package main

import (
	"fmt"
	"time"
)

/*
Games is an interface that defines the interactions needed to get game data
*/
type Games interface {
	GamesForDates(dates []time.Time)
	SetBaseURL(url string)
}

/*
LocalGames is an implementation of the Games interface
*/
type LocalGames struct {
	baseURL string
	cache   GameCache
	url     GameURL
}

/*
SetBaseURL sets the url for the baseball website
*/
func (lc *LocalGames) SetBaseURL(url string) {
	lc.baseURL = url
	if lc.cache == nil {
		lc.cache = &FSCache{}
	}
	lc.cache.SetBaseURL(url)
	if lc.url == nil {
		lc.url = &LocalURL{}
	}
	lc.url.SetBaseURL(url)
}

/*
GamesForDates returns all of the games associated with the dates
*/
func (lc *LocalGames) GamesForDates(dates []time.Time) {
	urls := lc.url.GetURLsForDates(dates)
	for _, url := range urls {
		fmt.Println(url.Path)
	}
}
