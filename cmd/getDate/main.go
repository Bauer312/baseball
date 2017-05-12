package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/Bauer312/baseball/pkg/dateslice"
	"github.com/Bauer312/baseball/pkg/util"
)

func main() {
	begDt := flag.String("beg", "", "Retrieve all games starting on this date.  Dates are in YYYYMMDD format")
	endDt := flag.String("end", "", "Retrieve all games ending on this date.  Dates are in YYYYMMDD format")
	dateString := flag.String("date", "", "Retrieve all games using text such as today, yesterday, thisweek, lastweek")

	flag.Parse()

	ds := dateslice.DateObjectsToSlice(*dateString, *begDt, *endDt)

	if ds != nil {
		for i, d := range ds {
			util.SetRoot("http://gd2.mlb.com/components/game/mlb", "/usr/local/share/baseball")
			dateURL, err := util.DateToURL(d)
			if err != nil {
				fmt.Println(err)
			}
			dateFS, err := util.URLToFSPath(dateURL)
			if err != nil {
				fmt.Println(err)
			}

			// Be kind to the web server, if there are multiple requests, wait 5 seconds between them
			if i > 0 {
				time.Sleep(5 * time.Second)
			}

			err = util.SaveURLToPath(dateURL, dateFS)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
