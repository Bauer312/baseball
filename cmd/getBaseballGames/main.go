/*
	Copyright 2019 Brian Bauer

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
	"time"

	"github.com/bauer312/baseball/pkg/dateslice"
)

func main() {
	date := flag.String("date", "yesterday", "Retreive data for a specific date (default is yesterday)")
	start := flag.String("start", "", "Retreive data for a date range (YYYYMMDD)")
	end := flag.String("end", "", "Retreive data for a date range (YYYYMMDD)")
	output := flag.String("output", "~/baseball", "Output location for downloaded files")

	flag.Parse()

	fmt.Printf("Raw data will be saved in %s\n", *output)

	var dates []time.Time
	if len(*start) > 0 {
		dates = dateslice.DateObjectsToSlice("", *start, *end)
	} else {
		dates = dateslice.DateStringToSlice(*date)
		if len(dates) == 0 {
			dates = dateslice.DateObjectsToSlice("", *date, *date)
		}
	}

	for i, dt := range dates {
		fmt.Printf("Downloading data for [%d] %s\n", i+1, dt.Format("20060102"))
	}
}
