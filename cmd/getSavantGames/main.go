package main

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

func main() {
	date := flag.String("date", "yesterday", "Retreive data for a specific date (default is yesterday)")
	start := flag.String("start", "", "Retreive data for a date range (YYYYMMDD)")
	end := flag.String("end", "", "Retreive data for a date range (YYYYMMDD)")
	output := flag.String("output", "", "Output location for downloaded files")
	url := flag.String("url", "https://baseballsavant.mlb.com", "Source location of data to download")

	flag.Parse()

	var dates []time.Time
	if len(*start) > 0 {
		dates = dateslice.DateObjectsToSlice("", *start, *end)
	} else {
		dates = dateslice.DateStringToSlice(*date)
		if len(dates) == 0 {
			dates = dateslice.DateObjectsToSlice("", *date, *date)
		}
	}

	fullOutputPath := validateOutput(*output)

	for i, dt := range dates {
		targetURL := dateToPath(*url, dt)
		fmt.Printf("Downloading data for [%d] %s (%s)\n", i+1, dt.Format("20060102"), targetURL)
		downloadFile(targetURL, filepath.Join(fullOutputPath, dt.Format("20060102")+".csv"))
	}
}

func dateToPath(baseURL string, date time.Time) string {
	year := date.Year()
	month := date.Month()
	day := date.Day()
	return fmt.Sprintf("%s/statcast_search/csv?all=true&hfPT=&hfAB=&hfBBT=&hfPR=&hfZ=&stadium=&hfBBL=&hfNewZones=&hfGT=R|&hfC=&hfSea=%04d|&hfSit=&player_type=pitcher&hfOuts=&opponent=&pitcher_throws=&batter_stands=&hfSA=&game_date_gt=%04d-%02d-%02d&game_date_lt=%04d-%02d-%02d&hfInfield=&team=&position=&hfOutfield=&hfRO=&home_road=&hfFlag=&hfPull=&metric_1=&hfInn=&min_pitches=0&min_results=0&group_by=name&sort_col=pitches&player_event_sort=h_launch_speed&sort_order=desc&min_pas=0&type=details&", baseURL, year, year, month, day, year, month, day)
}

func downloadFile(url, target string) {
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
