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

package util

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
)

/*
GameXMLStadium decribes the stadium structure present in the game.xml file
*/
type GameXMLStadium struct {
	ID            string `xml:"id,attr"`
	Name          string `xml:"name,attr"`
	VenueWhoKnows string `xml:"venue_w_chan_loc,attr"`
	Location      string `xml:"location,attr"`
}

/*
GameXMLTeam describes the team structure present in the game.xml file
*/
type GameXMLTeam struct {
	Type         string `xml:"type,attr"`
	Code         string `xml:"code,attr"`
	FileCode     string `xml:"file_code,attr"`
	Abbreviation string `xml:"abbrev,attr"`
	ID           string `xml:"id,attr"`
	Name         string `xml:"name,attr"`
	FullName     string `xml:"name_full,attr"`
	BriefName    string `xml:"name_brief,attr"`
	Wins         string `xml:"w,attr"`
	Losses       string `xml:"l,attr"`
	DivisionID   string `xml:"division_id,attr"`
	LeagueID     string `xml:"league_id,attr"`
	League       string `xml:"league,attr"`
}

/*
GameXMLGame describes the game structure present in the game.xml file
*/
type GameXMLGame struct {
	Type            string         `xml:"type,attr"`
	LocalGameTime   string         `xml:"local_game_time,attr"`
	GamePK          string         `xml:"game_pk,attr"`
	GameTimeEastern string         `xml:"game_time_et,attr"`
	GamedaySW       string         `xml:"gameday_sw,attr"`
	Teams           []GameXMLTeam  `xml:"team"`
	Stadium         GameXMLStadium `xml:"stadium"`
}

/*
ParseGameXML is a method that opens the locally-saved game.xml file and parses the
	contents into data structures.
*/
func ParseGameXML(path, date string, filePtr *os.File) error {
	gameID := filepath.Base(filepath.Dir(path))
	fileReader, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fileReader.Close()

	var g GameXMLGame
	decoder := xml.NewDecoder(fileReader)
	err = decoder.Decode(&g)
	if err != nil {
		return err
	}
	for _, team := range g.Teams {
		teamString := fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s\n",
			date,
			gameID,
			g.Type,
			g.LocalGameTime,
			g.GamePK,
			g.GameTimeEastern,
			g.GamedaySW,
			team.Type,
			team.Code,
			team.FileCode,
			team.Abbreviation,
			team.ID,
			team.Name,
			team.FullName,
			team.BriefName,
			team.Wins,
			team.Losses,
			team.DivisionID,
			team.LeagueID,
			team.League,
			g.Stadium.ID,
			g.Stadium.Name,
			g.Stadium.VenueWhoKnows,
			g.Stadium.Location,
		)
		_, err := filePtr.WriteString(teamString)
		if err != nil {
			fmt.Println(err)
		}
	}

	return nil
}
