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
ParseGameEventsXML is a method that opens the locally-saved game_events.xml file and parses the
	contents into data structures.
*/
func ParseGameEventsXML(path, date string, filePtr *os.File) error {
	gameID := filepath.Base(filepath.Dir(path))
	fileReader, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fileReader.Close()

	var g GameEventsXMLGame
	decoder := xml.NewDecoder(fileReader)
	err = decoder.Decode(&g)
	if err != nil {
		return err
	}

	initialPrefix := date + "|" + gameID

	for _, inning := range g.Innings {
		inningPrefix := initialPrefix + "|" + inning.Number

		topHalfPrefix := inningPrefix + "|" + "Top"
		printHalfInning(inning.TopHalf, topHalfPrefix, filePtr)
		bottomHalfPrefix := inningPrefix + "|" + "Bottom"
		printHalfInning(inning.BottomHalf, bottomHalfPrefix, filePtr)
	}

	return nil
}

func printHalfInning(innings GameEventsXMLHalfInning, linePrefix string, filePtr *os.File) {
	for _, atbat := range innings.AtBats {
		atBatPrefix := fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s",
			linePrefix,
			atbat.BatterNumber,
			atbat.Balls,
			atbat.Strikes,
			atbat.Outs,
			atbat.StartTFS,
			atbat.StartTFSZulu,
			atbat.Batter,
			atbat.Pitcher,
			atbat.EnglishDescription,
			atbat.EspañolDescription,
			atbat.EventNumber,
			atbat.EnglishEvent,
			atbat.EspañolEvent,
			atbat.PlayGUID,
			atbat.Score,
			atbat.HomeTeamRuns,
			atbat.AwayTeamRuns,
			atbat.FirstBasePlayer,
			atbat.SecondBasePlayer,
			atbat.ThirdBasePlayer,
		)
		for _, pitch := range atbat.Pitches {
			outputString := fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s\n",
				atBatPrefix,
				pitch.SVID,
				pitch.EnglishDescription,
				pitch.EspañolDescription,
				pitch.Type,
				pitch.StartSpeed,
				pitch.PitchType,
			)
			_, err := filePtr.WriteString(outputString)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	for _, action := range innings.Actions {
		outputString := fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s\n",
			linePrefix,
			action.Balls,
			action.Strikes,
			action.Outs,
			action.EnglishDescription,
			action.EspañolDescription,
			action.EnglishEvent,
			action.EspañolEvent,
			action.TFS,
			action.TFSZulu,
			action.Player,
			action.Pitch,
			action.EventNumber,
			action.HomeTeamRuns,
			action.AwayTeamRuns,
		)
		_, err := filePtr.WriteString(outputString)
		if err != nil {
			fmt.Println(err)
		}
	}
}
