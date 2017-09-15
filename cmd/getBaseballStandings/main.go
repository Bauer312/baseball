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

package main

import (
	"flag"
	"fmt"

	"github.com/bauer312/baseball/pkg/dateslice"
	"github.com/bauer312/baseball/pkg/reports"
	"github.com/bauer312/baseball/pkg/util"
)

func main() {
	reportDate := flag.String("date", "", "Retrieve a standings report on this date.  Dates are in YYYYMMDD format")
	output := flag.String("out", "", "Specify where to write data.  The default is the screen")

	flag.Parse()

	if len(*reportDate) == 0 {
		tomorrow := dateslice.Tomorrow()
		if len(tomorrow) == 1 {
			*reportDate = tomorrow[0].UTC().Format("2006-01-02")
		} else {
			fmt.Println("Unable to get a date from the system.")
			return
		}
	}
	fmt.Printf("Using %s as the date\n", *reportDate)

	db, err := util.GetDBConnection()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	reports.GetStandingsReport(db, *reportDate, "American League", "East", *output)
	reports.GetStandingsReport(db, *reportDate, "American League", "Central", *output)
	reports.GetStandingsReport(db, *reportDate, "American League", "West", *output)
	reports.GetStandingsReport(db, *reportDate, "National League", "East", *output)
	reports.GetStandingsReport(db, *reportDate, "National League", "Central", *output)
	reports.GetStandingsReport(db, *reportDate, "National League", "West", *output)
}
