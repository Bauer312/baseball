package main

/*
GameCache is an interface to a caching mechanism for game data.
	It is a good idea to be a nice citizen and avoid hitting the MLB servers too often.
*/
type GameCache interface {
	SetBaseURL(url string)
	GetURL(url string)
	InvalidateUrl(url string)
}

/*
FSCache is a basic filesystem cache that uses an environment variable to specify the location.
*/
type FSCache struct {
	baseURL string
}

/*
SetBaseURL provides the base URL that should be used as a starting point for the cache file names
*/
func (fsc *FSCache) SetBaseURL(url string) {
	fsc.baseURL = url
}

/*
GetURL will return the contents of the specified URL.
	If it exists in the cache, that will be returned. If not, it will be retrieved from the server.
*/
func (fsc *FSCache) GetURL(url string) {

}
