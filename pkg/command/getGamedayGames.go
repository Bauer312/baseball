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

package command

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

/*
GetGamedayGames contains information used to get the gameday
	data from the MLB website
*/
type GetGamedayGames struct {
	date   string
	start  string
	end    string
	output string
	url    string
}

/*
SetFlags creates the flags that are needed for this functionality
*/
func (ggg *GetGamedayGames) SetFlags(fs *flag.FlagSet, cmdMap map[string]*string) {
	cmdMap["date"] = fs.String("date", "yesterday", "Retreive data for a specific date (default is yesterday)")
	cmdMap["start"] = fs.String("start", "", "Retreive data for a date range (YYYYMMDD)")
	cmdMap["end"] = fs.String("end", "", "Retreive data for a date range (YYYYMMDD)")
	cmdMap["output"] = fs.String("output", "", "Output location for downloaded files")
	cmdMap["url"] = fs.String("url", "http://gd2.mlb.com", "Source location of data to download")

}

/*
Execute runs the functionality that produces the data needed
*/
func (ggg *GetGamedayGames) Execute(cmdMap map[string]*string) {
	ggg.date = *cmdMap["date"]
	ggg.start = *cmdMap["start"]
	ggg.end = *cmdMap["end"]
	ggg.output = *cmdMap["output"]
	ggg.url = *cmdMap["url"]

	var dates []time.Time
	if len(ggg.start) > 0 {
		dates = dateslice.DateObjectsToSlice("", ggg.start, ggg.end)
	} else {
		dates = dateslice.DateStringToSlice(ggg.date)
		if len(dates) == 0 {
			dates = dateslice.DateObjectsToSlice("", ggg.date, ggg.date)
		}
	}

	client := http.Client{Timeout: (45 * time.Second)}
	var datePaths datepath.DatePath

	datePaths.Init()
	go datePaths.ChannelListener(&client)

	var wg sync.WaitGroup
	wg.Add(1)

	go printFilePath(&wg, datePaths.FilePath, ggg.output)

	for i, dt := range dates {
		fmt.Printf("Downloading data for [%d] %s (%s)\n",
			i+1, dt.Format("20060102"), dateToPath(ggg.url, dt))
		datePaths.DatePath <- dateToPath(ggg.url, dt)
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
	client := http.Client{Timeout: (45 * time.Second)}
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
