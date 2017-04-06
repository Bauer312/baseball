package main

import "fmt"

func main() {
	cache := FSCache{}
	cache.SetBaseURL("http://gd2.mlb.com/components/game/mlb/")

	fmt.Println("Games")
}
