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

package pipelineStage

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

/*
GameFile contains the elements of the stage
*/
type GameFile struct {
	DataInput  chan string
	DataOutput chan FileOutputParameters
	wg         sync.WaitGroup
	rwg        sync.WaitGroup
}

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
ChannelListener should be run in a goroutine and will receive data on the input channel
	and the input control channel.  The parameters are converted to a slice of time
	elements and these elements are then sent out over the output channel.
*/
func (gF *GameFile) ChannelListener(client *http.Client) {
	for inputData := range gF.DataInput {
		gF.rwg.Add(1)
		resp, err := client.Get(inputData)
		if err != nil {
			fmt.Println(err.Error())
		}
		gF.tokenize(inputData, resp)
	}
	gF.rwg.Wait()

	//Tell the pipeline we are done
	gF.wg.Done()
}

func (gF *GameFile) tokenize(dataPath string, resp *http.Response) {
	defer resp.Body.Close()

	gameDate := gF.getDate(dataPath)
	dateElement, err := time.Parse("20060102", gameDate)
	if err != nil {
		gF.rwg.Done()
		return
	}

	var g GameXMLGame
	decoder := xml.NewDecoder(resp.Body)
	err = decoder.Decode(&g)
	if err != nil {
		gF.rwg.Done()
		return
	}

	gameString := fmt.Sprintf("%s|%s|%s|%s|%s|%s\n",
		gameDate,
		g.GamePK,
		g.Type,
		g.LocalGameTime,
		g.GameTimeEastern,
		g.GamedaySW,
	)
	gameOutput := FileOutputParameters{
		FileName:   "gameInfo",
		RecordDate: dateElement,
		DataRecord: gameString,
	}
	gF.DataOutput <- gameOutput

	for _, team := range g.Teams {
		teamString := fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s\n",
			gameDate,
			g.GamePK,
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
		)

		teamOutput := FileOutputParameters{
			FileName:   "teamInfo",
			RecordDate: dateElement,
			DataRecord: teamString,
		}
		gF.DataOutput <- teamOutput
	}

	stadiumString := fmt.Sprintf("%s|%s|%s|%s|%s|%s\n",
		gameDate,
		g.GamePK,
		g.Stadium.ID,
		g.Stadium.Name,
		g.Stadium.VenueWhoKnows,
		g.Stadium.Location,
	)

	stadiumOutput := FileOutputParameters{
		FileName:   "stadiumInfo",
		RecordDate: dateElement,
		DataRecord: stadiumString,
	}
	gF.DataOutput <- stadiumOutput

	gF.rwg.Done()
}

/*
Init will create all channels and other initialization needs.
	The DataInput channel is the output of any previous
	pipeline stage so it shouldn't be created here
*/
func (gF *GameFile) Init() error {
	gF.wg.Add(1)
	gF.DataOutput = make(chan FileOutputParameters, 5)

	return nil
}

/*
Stop will close the input channel, causing the Channel Listener to stop
*/
func (gF *GameFile) Stop() {
	close(gF.DataInput)
	gF.wg.Wait()
}

/*
getDate will attempt to discover the date from a URL
*/
func (gF *GameFile) getDate(dataPath string) string {
	year := 0
	month := 0
	day := 0
	components := strings.Split(dataPath, "/")
	for _, component := range components {
		if strings.HasPrefix(component, "year_") {
			year, _ = strconv.Atoi(component[5:])
		}
		if strings.HasPrefix(component, "month_") {
			month, _ = strconv.Atoi(component[6:])
		}
		if strings.HasPrefix(component, "day_") {
			day, _ = strconv.Atoi(component[4:])
		}
	}
	return fmt.Sprintf("%4d%02d%02d", year, month, day)
}
