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
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	records "github.com/bauer312/baseball/pkg/records"
)

/*
ScoreBoardFile contains the elements of the stage
*/
type ScoreBoardFile struct {
	DataInput      chan string
	DataOutput     chan string
	GameFileOutout chan string
	BaseURL        string
	Client         *http.Client
	wg             sync.WaitGroup
	rwg            sync.WaitGroup
}

/*
ScoreboardXMLGames describes the games structure present in the master_scoreboard.xml file
*/
type ScoreboardXMLGames struct {
	Year         int                 `xml:"year,attr"`
	Month        int                 `xml:"month,attr"`
	Day          int                 `xml:"day,attr"`
	LastModified string              `xml:"modified_date,attr"`
	NextDay      string              `xml:"next_day_date,attr"`
	Games        []ScoreboardXMLGame `xml:"game"`
}

/*
ScoreboardXMLGame describes the game structure present in the master_scoreboard.xml file
*/
type ScoreboardXMLGame struct {
	ID                    string                     `xml:"id,attr"`
	Venue                 string                     `xml:"venue,attr"`
	PK                    int                        `xml:"game_pk,attr"`
	Time                  string                     `xml:"time,attr"`
	DateTime              string                     `xml:"time_date,attr"`
	TimeDateAwLg          string                     `xml:"time_date_aw_lg,attr"`
	TimeDateHmLg          string                     `xml:"time_date_hm_lg,attr"`
	TimeZone              string                     `xml:"time_zone,attr"`
	AMPM                  string                     `xml:"ampm,attr"`
	FirstPitchET          string                     `xml:"first_pitch_et,attr"`
	AwayTime              string                     `xml:"away_time,attr"`
	AwayTimeZone          string                     `xml:"away_time_zone,attr"`
	AwayAMPM              string                     `xml:"away_ampm,attr"`
	HomeTime              string                     `xml:"home_time,attr"`
	HomeTimeZone          string                     `xml:"home_time_zone,attr"`
	HomeAMPM              string                     `xml:"home_ampm,attr"`
	GameType              string                     `xml:"game_type,attr"`
	TieBreakerSW          string                     `xml:"tiebreaker_sw,attr"`
	ResumeDate            string                     `xml:"resume_date,attr"`
	OriginalDate          string                     `xml:"original_date,attr"`
	TimeZoneAwLg          string                     `xml:"time_zone_aw_lg,attr"`
	TimeZoneHmLg          string                     `xml:"time_zone_hm_lg,attr"`
	TimeAwayLg            string                     `xml:"time_aw_lg,attr"`
	AwLgAMPM              string                     `xml:"aw_lg_ampm,attr"`
	TzAwLgGen             string                     `xml:"tz_aw_lg_gen,attr"`
	TimeHmLg              string                     `xml:"time_hm_lg,attr"`
	HmLgGen               string                     `xml:"hm_lg_gen,attr"`
	TzHmLgGen             string                     `xml:"tz_hm_lg_gen,attr"`
	VenueID               string                     `xml:"venue_id,attr"`
	ScheduledInnings      int                        `xml:"scheduled_innings,attr"`
	Description           string                     `xml:"description,attr"`
	AwayNameAbbrev        string                     `xml:"away_name_abbrev,attr"`
	HomeNameAbbrev        string                     `xml:"home_name_abbrev,attr"`
	AwayCode              string                     `xml:"away_code,attr"`
	AwayFileCode          string                     `xml:"away_file_code,attr"`
	AwayTeamID            string                     `xml:"away_team_id,attr"`
	AwayTeamCity          string                     `xml:"away_team_city,attr"`
	AwayTeamName          string                     `xml:"away_team_name,attr"`
	AwayDivision          string                     `xml:"away_division,attr"`
	AwayLeagueID          string                     `xml:"away_league_id,attr"`
	AwaySportCode         string                     `xml:"away_sport_code,attr"`
	HomeCode              string                     `xml:"home_code,attr"`
	HomeFileCode          string                     `xml:"home_file_code,attr"`
	HomeTeamID            string                     `xml:"home_team_id,attr"`
	HomeTeamCity          string                     `xml:"home_team_city,attr"`
	HomeTeamName          string                     `xml:"home_team_name,attr"`
	HomeDivision          string                     `xml:"home_division,attr"`
	HomeLeagueID          string                     `xml:"home_league_id,attr"`
	HomeSportCode         string                     `xml:"home_sport_code,attr"`
	Day                   string                     `xml:"day,attr"`
	GamedaySW             string                     `xml:"gameday_sw,attr"`
	DoubleHeaderSW        string                     `xml:"double_header_sw,attr"`
	GameNumber            int                        `xml:"game_nbr,attr"`
	TBDFlag               string                     `xml:"tbd_flag,attr"`
	AwayGamesBack         string                     `xml:"away_games_back,attr"`
	HomeGamesBack         string                     `xml:"home_games_back,attr"`
	AwayGamesBackWildcard string                     `xml:"away_games_back_wildcard,attr"`
	HomeGamesBackWildcard string                     `xml:"home_games_back_wildcard,attr"`
	VenueWChanLoc         string                     `xml:"venue_w_chan_loc,attr"`
	Location              string                     `xml:"location,attr"`
	GameDay               string                     `xml:"gameday,attr"`
	AwayWins              int                        `xml:"away_win,attr"`
	AwayLosses            int                        `xml:"away_loss,attr"`
	HomeWins              int                        `xml:"home_win,attr"`
	HomeLosses            int                        `xml:"home_loss,attr"`
	GameDataDirectory     string                     `xml:"game_data_directory,attr"`
	League                string                     `xml:"league,attr"`
	Status                ScoreboardXMLGameStatus    `xml:"status"`
	Linescore             ScoreboardXMLGameLinescore `xml:"linescore"`
}

/*
ScoreboardXMLGameStatus describes the game structure present in the master_scoreboard.xml file
*/
type ScoreboardXMLGameStatus struct {
	Status      string `xml:"status,attr"`
	Ind         string `xml:"ind,attr"`
	Reason      string `xml:"reason,attr"`
	Inning      int    `xml:"inning,attr"`
	TopInning   string `xml:"top_inning,attr"`
	Balls       int    `xml:"b,attr"`
	Strikes     int    `xml:"s,attr"`
	Outs        int    `xml:"o,attr"`
	InningState string `xml:"inning_state,attr"`
	Note        string `xml:"note,attr"`
	Perfect     string `xml:"is_perfect_game,attr"`
	NoHitter    string `xml:"is_no_hitter,attr"`
}

/*
ScoreboardXMLGameLinescore describes the inning by inning linescore structure present in the
	master_scoreboard.xml file
*/
type ScoreboardXMLGameLinescore struct {
	Innings []struct {
		Away int `xml:"away,attr"`
		Home int `xml:"home,attr"`
	} `xml:"inning"`
	Runs struct {
		Away int `xml:"away,attr"`
		Home int `xml:"home,attr"`
		Diff int `xml:"diff,attr"`
	} `xml:"r"`
	Hits struct {
		Away int `xml:"away,attr"`
		Home int `xml:"home,attr"`
	} `xml:"h"`
	Errors struct {
		Away int `xml:"away,attr"`
		Home int `xml:"home,attr"`
	} `xml:"e"`
	HR struct {
		Away int `xml:"away,attr"`
		Home int `xml:"home,attr"`
	} `xml:"hr"`
	SB struct {
		Away int `xml:"away,attr"`
		Home int `xml:"home,attr"`
	} `xml:"sb"`
	SO struct {
		Away int `xml:"away,attr"`
		Home int `xml:"home,attr"`
	} `xml:"so"`
}

/*
Run should be run in a goroutine and will receive URLs on the input channel.
    It will add the scoreboard file (master_scoreboard.xml) to the URL, retrieve it,
    and then parse it.  There will be two kinds of output:
        1. Game path and primary key to assist in getting further data about the game
        2. Game data contained in the file to be stored for future use
*/
func (sbF *ScoreBoardFile) Run() {
	for inputData := range sbF.DataInput {
		if strings.HasSuffix(inputData, "/") == false {
			inputData = inputData + "/"
		}
		inputData = inputData + "master_scoreboard.xml"

		sbF.rwg.Add(1)
		resp, err := sbF.Client.Get(inputData)
		if err != nil {
			fmt.Println(err.Error())
		}
		sbF.tokenize(inputData, resp)
	}
	sbF.rwg.Wait()

	//Tell the pipeline we are done
	sbF.wg.Done()
}

/*
Init will create all channels and other initialization needs.
	The DataInput channel is the output of any previous
	pipeline stage so it shouldn't be created here
*/
func (sbF *ScoreBoardFile) Init() error {
	sbF.wg.Add(1)
	sbF.DataOutput = make(chan string)
	sbF.GameFileOutout = make(chan string)

	return nil
}

/*
Stop will close the input channel, causing the Channel Listener to stop
*/
func (sbF *ScoreBoardFile) Stop() {
	close(sbF.DataInput)
	sbF.wg.Wait()
}

func (sbF *ScoreBoardFile) tokenize(dataPath string, resp *http.Response) {
	defer resp.Body.Close()
	defer sbF.rwg.Done()

	var sb ScoreboardXMLGames
	decoder := xml.NewDecoder(resp.Body)
	err := decoder.Decode(&sb)
	if err != nil {
		return
	}

	EasternLocation, err := time.LoadLocation("America/New_York")
	if err != nil {
		fmt.Printf("Unable to get location for the eastern time zone: %s", err.Error())
		return
	}
	timeFormat := "2006/01/02 3:04PM"
	//t, _ := time.ParseInLocation(longForm, "Jul 9, 2012 at 5:02am (CEST)", loc)
	//fmt.Println(t)

	venues := make([]records.VenueRecord, len(sb.Games))
	leagues := make([]records.LeagueRecord, len(sb.Games)*2)
	divisions := make([]records.DivisionRecord, len(sb.Games)*2)
	teams := make([]records.TeamRecord, len(sb.Games)*2)
	standings := make([]records.StandingRecord, len(sb.Games)*2)
	games := make([]records.GameRecord, len(sb.Games))
	gameStatuses := make([]records.GameStatusRecord, len(sb.Games))

	for i, game := range sb.Games {
		//First, output the game directory
		sbF.GameFileOutout <- game.GameDataDirectory

		//Second, create all data records
		timeString := game.DateTime + game.AMPM
		var usedLocation *time.Location
		switch game.TimeZone {
		case "ET":
			usedLocation = EasternLocation
		default:
			fmt.Printf("Unexpected Time Zone encountered: %s", game.TimeZone)
			return
		}
		gameTime, err := time.ParseInLocation(timeFormat, timeString, usedLocation)
		if err != nil {
			fmt.Printf("Unable to parse the timestamp of game %d (%s): %s", i, timeString, err.Error())
			return
		}
		venueID, err := strconv.ParseInt(game.VenueID, 10, 64)
		if err != nil {
			fmt.Printf("Unable to parse the venue ID of game %d (%s): %s", i, game.VenueID, err.Error())
			return
		}
		awayLeagueID, err := strconv.ParseInt(game.AwayLeagueID, 10, 64)
		if err != nil {
			fmt.Printf("Unable to parse the away league ID of game %d (%s): %s", i, game.AwayLeagueID, err.Error())
			return
		}
		homeLeagueID, err := strconv.ParseInt(game.HomeLeagueID, 10, 64)
		if err != nil {
			fmt.Printf("Unable to parse the home league ID of game %d (%s): %s", i, game.HomeLeagueID, err.Error())
			return
		}
		var awayLeagueName string
		switch awayLeagueID {
		case 103:
			awayLeagueName = "American League"
		case 104:
			awayLeagueName = "National League"
		default:
			fmt.Printf("Unexpected league ID ecountered: %d", awayLeagueID)
			return
		}
		var homeLeagueName string
		switch homeLeagueID {
		case 103:
			homeLeagueName = "American League"
		case 104:
			homeLeagueName = "National League"
		default:
			fmt.Printf("Unexpected league ID ecountered: %d", homeLeagueID)
			return
		}

		var awayDivisionName string
		switch game.AwayDivision {
		case "E":
			awayDivisionName = "East"
		case "C":
			awayDivisionName = "Central"
		case "W":
			awayDivisionName = "West"
		default:
			fmt.Printf("Unexpected division code ecountered: %s", game.AwayDivision)
			return
		}

		var homeDivisionName string
		switch game.HomeDivision {
		case "E":
			homeDivisionName = "East"
		case "C":
			homeDivisionName = "Central"
		case "W":
			homeDivisionName = "West"
		default:
			fmt.Printf("Unexpected division code ecountered: %s", game.HomeDivision)
			return
		}

		awayTeamID, err := strconv.ParseInt(game.AwayTeamID, 10, 64)
		if err != nil {
			fmt.Printf("Unable to parse the away team ID of game %d (%s): %s", i, game.AwayTeamID, err.Error())
			return
		}
		homeTeamID, err := strconv.ParseInt(game.HomeTeamID, 10, 64)
		if err != nil {
			fmt.Printf("Unable to parse the home team ID of game %d (%s): %s", i, game.HomeTeamID, err.Error())
			return
		}

		var topOfInning bool
		switch game.Status.TopInning {
		case "Y":
			topOfInning = true
		case "N":
			topOfInning = false
		default:
			fmt.Printf("Unexpected top of inning code encountered: %s", game.Status.TopInning)
			return
		}

		var perfectGame bool
		switch game.Status.Perfect {
		case "Y":
			perfectGame = true
		case "N":
			perfectGame = false
		default:
			fmt.Printf("Unexpected perfect game code encountered: %s", game.Status.Perfect)
			return
		}

		var noHitter bool
		switch game.Status.NoHitter {
		case "Y":
			noHitter = true
		case "N":
			noHitter = false
		default:
			fmt.Printf("Unexpected no hitter code encountered: %s", game.Status.NoHitter)
			return
		}

		venues[i] = records.VenueRecord{
			RecordName:    "VenueRecord",
			EffectiveDate: gameTime,
			ID:            venueID,
			Name:          game.Venue,
			Location:      game.Location,
			Channel:       game.VenueWChanLoc,
		}
		leagues[i*2] = records.LeagueRecord{
			RecordName:    "LeagueRecord",
			EffectiveDate: gameTime,
			ID:            awayLeagueID,
			Name:          awayLeagueName,
			SportCode:     game.AwaySportCode,
		}
		leagues[i*2+1] = records.LeagueRecord{
			RecordName:    "LeagueRecord",
			EffectiveDate: gameTime,
			ID:            homeLeagueID,
			Name:          homeLeagueName,
			SportCode:     game.HomeSportCode,
		}
		divisions[i*2] = records.DivisionRecord{
			RecordName:    "DivisionRecord",
			EffectiveDate: gameTime,
			Name:          awayDivisionName,
			Code:          game.AwayDivision,
		}
		divisions[i*2+1] = records.DivisionRecord{
			RecordName:    "DivisionRecord",
			EffectiveDate: gameTime,
			Name:          homeDivisionName,
			Code:          game.HomeDivision,
		}
		teams[i*2] = records.TeamRecord{
			RecordName:    "TeamRecord",
			EffectiveDate: gameTime,
			ID:            awayTeamID,
			Name:          game.AwayTeamName,
			Code:          game.AwayCode,
			City:          game.AwayTeamCity,
			LeagueID:      awayLeagueID,
			Division:      game.AwayDivision,
		}
		teams[i*2+1] = records.TeamRecord{
			RecordName:    "TeamRecord",
			EffectiveDate: gameTime,
			ID:            homeTeamID,
			Name:          game.HomeTeamName,
			Code:          game.HomeCode,
			City:          game.HomeTeamCity,
			LeagueID:      homeLeagueID,
			Division:      game.HomeDivision,
		}
		standings[i*2] = records.StandingRecord{
			RecordName:        "StandingRecord",
			EffectiveDate:     gameTime,
			TeamID:            awayTeamID,
			Wins:              game.AwayWins,
			Losses:            game.AwayLosses,
			GamesBack:         game.AwayGamesBack,
			WildcardGamesBack: game.AwayGamesBackWildcard,
		}
		standings[i*2+1] = records.StandingRecord{
			RecordName:        "StandingRecord",
			EffectiveDate:     gameTime,
			TeamID:            homeTeamID,
			Wins:              game.HomeWins,
			Losses:            game.HomeLosses,
			GamesBack:         game.HomeGamesBack,
			WildcardGamesBack: game.HomeGamesBackWildcard,
		}
		games[i] = records.GameRecord{
			RecordName:       "GameRecord",
			EffectiveDate:    gameTime,
			ID:               int64(game.PK),
			ResumeDate:       game.ResumeDate,
			OriginalDate:     game.OriginalDate,
			GameType:         game.GameType,
			Tiebreaker:       game.TieBreakerSW,
			GameDay:          game.GamedaySW,
			DoubleHeader:     game.DoubleHeaderSW,
			GameNumber:       game.GameNumber,
			TBDFlag:          game.TBDFlag,
			Interleague:      game.League,
			ScheduledInnings: game.ScheduledInnings,
			Description:      game.Description,
			VenueID:          venueID,
			AwayTeamID:       awayTeamID,
			HomeTeamID:       homeTeamID,
		}

		gameStatuses[i] = records.GameStatusRecord{
			RecordName:     "GameStatusRecord",
			EffectiveDate:  gameTime,
			ID:             int64(game.PK),
			Status:         game.Status.Status,
			Ind:            game.Status.Ind,
			Reason:         game.Status.Reason,
			CurrentInning:  game.Status.Inning,
			TopOfInning:    topOfInning,
			Balls:          game.Status.Balls,
			Strikes:        game.Status.Strikes,
			Outs:           game.Status.Outs,
			InningState:    game.Status.InningState,
			Note:           game.Status.Note,
			PerfectGame:    perfectGame,
			NoHitter:       noHitter,
			AwayTeamRuns:   game.Linescore.Runs.Away,
			HomeTeamRuns:   game.Linescore.Runs.Home,
			AwayTeamHits:   game.Linescore.Hits.Away,
			HomeTeamHits:   game.Linescore.Hits.Home,
			AwayTeamErrors: game.Linescore.Errors.Away,
			HomeTeamErrors: game.Linescore.Errors.Home,
			AwayTeamHR:     game.Linescore.HR.Away,
			HomeTeamHR:     game.Linescore.HR.Home,
			AwayTeamSB:     game.Linescore.SB.Away,
			HomeTeamSB:     game.Linescore.SB.Home,
			AwayTeamSO:     game.Linescore.SO.Away,
			HomeTeamSO:     game.Linescore.SO.Home,
		}

		gameStatuses[i].Innings = make([]records.InningScoreRecord, len(game.Linescore.Innings))
		for y, gameInning := range game.Linescore.Innings {
			gameStatuses[i].Innings[y] = records.InningScoreRecord{
				RecordName:    "InningScoreRecord",
				EffectiveDate: gameTime,
				GameID:        int64(game.PK),
				Inning:        y + 1,
				AwayTeamRuns:  gameInning.Away,
				HomeTeamRuns:  gameInning.Home,
			}
		}

		sbF.sendJSONToOutput(json.Marshal(venues[i]))
		sbF.sendJSONToOutput(json.Marshal(leagues[i*2]))
		sbF.sendJSONToOutput(json.Marshal(leagues[i*2+1]))
		sbF.sendJSONToOutput(json.Marshal(divisions[i*2]))
		sbF.sendJSONToOutput(json.Marshal(divisions[i*2+1]))
		sbF.sendJSONToOutput(json.Marshal(teams[i*2]))
		sbF.sendJSONToOutput(json.Marshal(teams[i*2+1]))
		sbF.sendJSONToOutput(json.Marshal(standings[i*2]))
		sbF.sendJSONToOutput(json.Marshal(standings[i*2+1]))
		sbF.sendJSONToOutput(json.Marshal(games[i]))
		sbF.sendJSONToOutput(json.Marshal(gameStatuses[i]))
		for _, isR := range gameStatuses[i].Innings {
			sbF.sendJSONToOutput(json.Marshal(isR))
		}
	}
}

func (sbF *ScoreBoardFile) sendJSONToOutput(rep []byte, err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
	sbF.DataOutput <- string(rep)
}
