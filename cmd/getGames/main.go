package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/net/html"

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
			fileReader, err := os.Open(dateFS)
			if err != nil {
				fmt.Println(err)
			}
			defer fileReader.Close()
			htmlTokenizer := html.NewTokenizer(fileReader)
			for {
				tt := htmlTokenizer.Next()
				if tt == html.ErrorToken {
					break
				}
				if tt == html.StartTagToken {
					t := htmlTokenizer.Token()

					isAnchor := t.Data == "a"
					if isAnchor {
						for _, a := range t.Attr {
							if a.Key == "href" {
								if strings.HasPrefix(a.Val, "gid_") {
									gameURLs, err := util.GameToURLs(a.Val)
									if err != nil {
										fmt.Println(err)
										break
									}
									for i, gameURL := range gameURLs {
										gameFS, err := util.URLToFSPath(gameURL)
										if err != nil {
											fmt.Println(err)
										}
										// Be kind to the web server, if there are multiple requests, wait 5 seconds between them
										if i > 0 {
											time.Sleep(5 * time.Second)
										}

										err = util.SaveURLToPath(gameURL, gameFS)
										if err != nil {
											fmt.Println(err)
										}
									}
								}

								break
							}
						}
					}
				}
			}
		}
	}
}
