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
	"io"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/bauer312/baseball/pkg/dateslice"
)

/*
GetSavantGames contains information used to get the Baseball Savant
	data from the Baseball Savant website
*/
type GetSavantGames struct {
	date   string
	start  string
	end    string
	output string
	url    string
}

/*
SetFlags creates the flags that are needed for this functionality
*/
func (gsg *GetSavantGames) SetFlags(fs *flag.FlagSet, cmdMap map[string]*string) {
	cmdMap["date"] = fs.String("date", "yesterday", "Retreive data for a specific date (default is yesterday)")
	cmdMap["start"] = fs.String("start", "", "Retreive data for a date range (YYYYMMDD)")
	cmdMap["end"] = fs.String("end", "", "Retreive data for a date range (YYYYMMDD)")
	cmdMap["output"] = fs.String("output", "", "Output location for downloaded files")
	cmdMap["url"] = fs.String("url", "https://baseballsavant.mlb.com", "Source location of data to download")

}

/*
Execute runs the functionality that produces the data needed
*/
func (gsg *GetSavantGames) Execute(cmdMap map[string]*string) {
	gsg.date = *cmdMap["date"]
	gsg.start = *cmdMap["start"]
	gsg.end = *cmdMap["end"]
	gsg.output = *cmdMap["output"]
	gsg.url = *cmdMap["url"]

	var dates []time.Time
	if len(gsg.start) > 0 {
		dates = dateslice.DateObjectsToSlice("", gsg.start, gsg.end)
	} else {
		dates = dateslice.DateStringToSlice(gsg.date)
		if len(dates) == 0 {
			dates = dateslice.DateObjectsToSlice("", gsg.date, gsg.date)
		}
	}

	fullOutputPath := validateOutput(gsg.output)

	for i, dt := range dates {
		targetURL := savantdateToPath(gsg.url, dt)
		fmt.Printf("Downloading data for [%d] %s (%s)\n", i+1, dt.Format("20060102"), targetURL)
		savantdownloadFile(targetURL, filepath.Join(fullOutputPath, dt.Format("20060102")+".csv"))
	}
}

func savantdateToPath(baseURL string, date time.Time) string {
	year := date.Year()
	month := date.Month()
	day := date.Day()
	return fmt.Sprintf("%s/statcast_search/csv?all=true&hfPT=&hfAB=&hfBBT=&hfPR=&hfZ=&stadium=&hfBBL=&hfNewZones=&hfGT=R|&hfC=&hfSea=%04d|&hfSit=&player_type=pitcher&hfOuts=&opponent=&pitcher_throws=&batter_stands=&hfSA=&game_date_gt=%04d-%02d-%02d&game_date_lt=%04d-%02d-%02d&hfInfield=&team=&position=&hfOutfield=&hfRO=&home_road=&hfFlag=&hfPull=&metric_1=&hfInn=&min_pitches=0&min_results=0&group_by=name&sort_col=pitches&player_event_sort=h_launch_speed&sort_order=desc&min_pas=0&type=details&", baseURL, year, year, month, day, year, month, day)
}

func savantdownloadFile(url, target string) {
	fmt.Printf("Target: %s\n", target)
	client := http.Client{Timeout: (45 * time.Second)}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		defer resp.Body.Close()
		f, err := os.OpenFile(target, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, resp.Body)
	}
}

func validateOutput(output string) string {

	if len(output) == 0 {
		usr, err := user.Current()
		if err != nil {
			fmt.Println("Unable to determine user storage location")
		}
		output = filepath.Join(usr.HomeDir, "baseball")
	}
	basePath := filepath.Join(output, "savant/")
	err := os.MkdirAll(basePath, 0740)
	if err != nil {
		fmt.Println("Unable to validate storage location")
	}
	fmt.Printf("Storage Location: %s\n", basePath)
	return basePath
}
