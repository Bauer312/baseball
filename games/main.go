package main

import (
	"flag"
	"fmt"
)

/*
	A typical use case for this program is to ask for the games that occurred:
		* On a specific day
			* Today
			* Yesterday
			* An arbitrary day
		* In a specific month
			* This month
			* Last month
			* An arbitrary month
		* In a specific year
			* This year
			* Last year
			* An arbitrary year

	The purpose is to get data that will be used to either load into a database
		or do some sort of analysis.
*/

func main() {
	datePtr := flag.String("date", "today", "Retrieve all games on this date.  Dates are in YYYYMMDD format")

	flag.Parse()

	ds := LastWeek()
	games := LocalGames{}
	games.SetBaseURL("http://gd2.mlb.com/components/game/mlb/")
	games.GamesForDates(ds)

	fmt.Printf("Games for %s\n", *datePtr)
}
