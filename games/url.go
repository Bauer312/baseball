package main

/*
GameURL is an interface to a URL-building mechanism for game data.
	The MLB website has a pretty nice URL structure; given some
	parameters a URL can be built so that resource can be retrieved.
*/
type GameURL interface {
	SetBaseURL(url string)
}
