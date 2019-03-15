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

package pipelinestage

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

/*
GameEventsFile contains the elements of the stage
*/
type GameEventsFile struct {
	DataInput  chan string
	DataOutput chan FileOutputParameters
	wg         sync.WaitGroup
	rwg        sync.WaitGroup
}

/*
GameEventsXMLPitch describes the pitch structure present in the game_events.xml file
*/
type GameEventsXMLPitch struct {
	SVID               string `xml:"sv_id,attr"`
	EnglishDescription string `xml:"des,attr"`
	EspañolDescription string `xml:"des_es,attr"`
	Type               string `xml:"type,attr"`
	StartSpeed         string `xml:"start_speed,attr"`
	PitchType          string `xml:"pitch_type,attr"`
}

/*
GameEventsXMLAtBat describes the at-bat structure present in the game_events.xml file
*/
type GameEventsXMLAtBat struct {
	BatterNumber       string               `xml:"num,attr"`
	Balls              string               `xml:"b,attr"`
	Strikes            string               `xml:"s,attr"`
	Outs               string               `xml:"o,attr"`
	StartTFS           string               `xml:"start_tfs,attr"`
	StartTFSZulu       string               `xml:"start_tfs_zulu,attr"`
	Batter             string               `xml:"batter,attr"`
	Pitcher            string               `xml:"pitcher,attr"`
	EnglishDescription string               `xml:"des,attr"`
	EspañolDescription string               `xml:"des_es,attr"`
	EventNumber        string               `xml:"event_num,attr"`
	EnglishEvent       string               `xml:"event,attr"`
	EspañolEvent       string               `xml:"event_es,attr"`
	PlayGUID           string               `xml:"play_guid,attr"`
	Score              string               `xml:"score,attr"`
	HomeTeamRuns       string               `xml:"home_team_runs,attr"`
	AwayTeamRuns       string               `xml:"away_team_runs,attr"`
	FirstBasePlayer    string               `xml:"b1,attr"`
	SecondBasePlayer   string               `xml:"b2,attr"`
	ThirdBasePlayer    string               `xml:"b3,attr"`
	Pitches            []GameEventsXMLPitch `xml:"pitch"`
}

/*
GameEventsXMLAction describes the action structure present in the game_events.xml file
*/
type GameEventsXMLAction struct {
	Balls              string `xml:"b,attr"`
	Strikes            string `xml:"s,attr"`
	Outs               string `xml:"o,attr"`
	EnglishDescription string `xml:"des,attr"`
	EspañolDescription string `xml:"des_es,attr"`
	EnglishEvent       string `xml:"event,attr"`
	EspañolEvent       string `xml:"event_es,attr"`
	TFS                string `xml:"tfs,attr"`
	TFSZulu            string `xml:"tfs_zulu,attr"`
	Player             string `xml:"player,attr"`
	Pitch              string `xml:"pitch,attr"`
	EventNumber        string `xml:"event_num,attr"`
	HomeTeamRuns       string `xml:"home_team_runs,attr"`
	AwayTeamRuns       string `xml:"away_team_runs,attr"`
}

/*
GameEventsXMLHalfInning describes the half-inning structure present in the game_events.xml file
*/
type GameEventsXMLHalfInning struct {
	AtBats  []GameEventsXMLAtBat  `xml:"atbat"`
	Actions []GameEventsXMLAction `xml:"action"`
}

/*
GameEventsXMLInning describes the inning structure present in the game_events.xml file
*/
type GameEventsXMLInning struct {
	Number     string                  `xml:"num,attr"`
	TopHalf    GameEventsXMLHalfInning `xml:"top"`
	BottomHalf GameEventsXMLHalfInning `xml:"bottom"`
}

/*
GameEventsXMLGame describes the game structure present in the game_events.xml file
*/
type GameEventsXMLGame struct {
	Innings []GameEventsXMLInning `xml:"inning"`
}

/*
ChannelListener should be run in a goroutine and will receive data on the input channel
	and the input control channel.  The parameters are converted to a slice of time
	elements and these elements are then sent out over the output channel.
*/
func (gE *GameEventsFile) ChannelListener(client *http.Client) {
	for inputData := range gE.DataInput {
		gE.rwg.Add(1)
		resp, err := client.Get(inputData)
		if err != nil {
			fmt.Println(err.Error())
		}
		gE.tokenize(inputData, resp)
	}
	gE.rwg.Wait()

	//Tell the pipeline we are done
	gE.wg.Done()
}

func (gE *GameEventsFile) tokenize(dataPath string, resp *http.Response) {
	defer resp.Body.Close()

	gameDate := gE.getDate(dataPath)

	var g GameXMLGame
	decoder := xml.NewDecoder(resp.Body)
	err := decoder.Decode(&g)
	if err != nil {
		gE.rwg.Done()
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
		DataRecord: gameString,
	}
	gE.DataOutput <- gameOutput

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
			DataRecord: teamString,
		}
		gE.DataOutput <- teamOutput
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
		DataRecord: stadiumString,
	}
	gE.DataOutput <- stadiumOutput

	gE.rwg.Done()
}

/*
Init will create all channels and other initialization needs.
	The DataInput channel is the output of any previous
	pipeline stage so it shouldn't be created here
*/
func (gE *GameEventsFile) Init() error {
	gE.wg.Add(1)
	gE.DataOutput = make(chan FileOutputParameters, 5)

	return nil
}

/*
Stop will close the input channel, causing the Channel Listener to stop
*/
func (gE *GameEventsFile) Stop() {
	close(gE.DataInput)
	gE.wg.Wait()
}

/*
getDate will attempt to discover the date from a URL
*/
func (gE *GameEventsFile) getDate(dataPath string) string {
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
