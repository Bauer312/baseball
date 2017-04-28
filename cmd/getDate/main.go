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

	var ds []time.Time

	if len(*dateString) != 0 {
		switch *dateString {
		case "today":
			ds = dateslice.Today()
		case "yesterday":
			ds = dateslice.Yesterday()
		case "thisweek":
			ds = dateslice.ThisWeek()
		case "lastweek":
			ds = dateslice.LastWeek()
		case "thismonth":
			ds = dateslice.ThisMonth()
		case "lastmonth":
			ds = dateslice.LastMonth()
		}
	}

	if len(*begDt) != 0 {
		if len(*endDt) != 0 {
			ds = dateslice.RangeString(*begDt, *endDt)
		} else {
			ds = dateslice.RangeString(*begDt, *begDt)
		}
	}

	if ds != nil {
		for _, d := range ds {
			util.SetRoot("http://gd2.mlb.com/components/game/mlb", "/usr/local/share/baseball")
			dateURL, err := util.DateToURL(d)
			if err != nil {
				fmt.Println(err)
			}
			dateFS, err := util.URLToFSPath(dateURL)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("Save %s to %s\n", dateURL, dateFS)
		}
	}
}
