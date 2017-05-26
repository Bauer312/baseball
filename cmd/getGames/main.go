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
		var rsrc util.Resource
		rsrc.Roots("http://gd2.mlb.com/components/game/mlb", "/usr/local/share/baseball")
		for _, d := range ds {
			tDefs, err := rsrc.Date(d)
			if err != nil {
				fmt.Println(err)
			}

			for _, tDef := range tDefs {
				fileReader, err := os.Open(tDef.Target)
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
										gDefs, err := rsrc.Game(a.Val)
										if err != nil {
											fmt.Println(err)
											break
										}
										for i, gDef := range gDefs {
											// Be kind to the web server, if there are multiple requests, wait 3 seconds between them
											if i > 0 {
												time.Sleep(3 * time.Second)
											}

											err = util.SaveURLToPath(gDef.Source, gDef.Target)
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
}
