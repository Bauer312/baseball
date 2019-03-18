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
	"net/http"
	"sync"
	"time"

	"github.com/bauer312/baseball/pkg/datepath"
	"github.com/bauer312/baseball/pkg/dateslice"
	"github.com/bauer312/baseball/pkg/filepath"
)

func main() {
	date := flag.String("date", "yesterday", "Retreive data for a specific date (default is yesterday)")
	start := flag.String("start", "", "Retreive data for a date range (YYYYMMDD)")
	end := flag.String("end", "", "Retreive data for a date range (YYYYMMDD)")
	output := flag.String("output", "~/baseball", "Output location for downloaded files")
	url := flag.String("url", "http://gd2.mlb.com", "Source location of data to download")

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

	client := http.Client{Timeout: (10 * time.Second)}
	var datePaths datepath.DatePath

	datePaths.Init()
	go datePaths.ChannelListener(&client)

	var wg sync.WaitGroup
	wg.Add(1)

	go printFilePath(&wg, datePaths.FilePath, *output)

	for i, dt := range dates {
		fmt.Printf("Downloading data for [%d] %s (%s)\n",
			i+1, dt.Format("20060102"), dateToPath(*url, dt))
		datePaths.DatePath <- dateToPath(*url, dt)
	}

	datePaths.Done()

	wg.Wait()
}

func dateToPath(baseURL string, date time.Time) string {
	year := date.Year()
	month := date.Month()
	day := date.Day()
	return fmt.Sprintf("%s/components/game/mlb/year_%04d/month_%02d/day_%02d", baseURL, year, month, day)
}

func printFilePath(wg *sync.WaitGroup, paths chan string, output string) {
	client := http.Client{Timeout: (10 * time.Second)}
	var filePaths filepath.FilePath
	filePaths.Init(output)
	go filePaths.ChannelListener(&client)
	for path := range paths {
		//fmt.Printf("\t%s\n", path)
		filePaths.FilePath <- path
	}
	filePaths.Done()
	wg.Done()
}
