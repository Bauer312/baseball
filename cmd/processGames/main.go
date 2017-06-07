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
	"path/filepath"

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
			dateString := d.Format("20060102")
			processedPath, err := util.DateToProcessedFileNoSideEffects(d, "http://gd2.mlb.com/components/game/mlb", "/usr/local/share/baseball", "game.dat")
			if err != nil {
				fmt.Println(err)
			}
			err = util.VerifyFSDirectory(processedPath)
			if err != nil {
				fmt.Println(err)
			}
			filePtr, err := os.Create(processedPath)
			if err != nil {
				fmt.Println(err)
			}
			defer filePtr.Close()

			fmt.Println(processedPath)
			tDefs, err := rsrc.Date(d)
			if err != nil {
				fmt.Println(err)
			}

			for _, tDef := range tDefs {
				gameIDs, err := util.ParseGameFile(tDef.Target)
				if err != nil {
					fmt.Println(err)
				}
				for _, gameID := range gameIDs {
					gDefs, err := rsrc.Game(gameID)
					if err != nil {
						fmt.Println(err)
						break
					}
					for _, gDef := range gDefs {
						fileName := filepath.Base(gDef.Target)
						switch fileName {
						case "game.xml":
							gameXMLContents, err := util.ParseGameXML(gDef.Target)
							if err != nil {
								fmt.Println(err)
							}
							for _, gameString := range gameXMLContents {
								fmt.Print(dateString + "|" + gameString)
								_, err = filePtr.WriteString(dateString + "|" + gameString)
								if err != nil {
									fmt.Println(err)
								}
							}
						case "game_events.xml":
							gameEventsXMLContents, err := util.ParseGameEventsXML(gDef.Target)
							if err != nil {
								fmt.Println(err)
							}
							for _, gameEventString := range gameEventsXMLContents {
								fmt.Print(gameEventString)
							}
						}
					}
				}
			}
		}
	}
}
