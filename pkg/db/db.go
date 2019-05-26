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

package db

import (
	"database/sql"
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/lib/pq"
)

/*
BaseballDB holds data that allows interaction with a database
*/
type BaseballDB struct {
	dbConn *sql.DB
}

/*
ConnectToPostgres does something
*/
func (bdb *BaseballDB) ConnectToPostgres(prefix string) error {
	pWord := os.Getenv(prefix + "_PASS")
	uName := os.Getenv(prefix + "_USER")
	dbName := os.Getenv(prefix + "_DBNAME")
	sslMode := os.Getenv(prefix + "_SSLMODE")
	cString := fmt.Sprintf("user=%s dbname=%s password='%s' sslmode=%s",
		uName, dbName, pWord, sslMode)
	db, err := sql.Open("postgres", cString)
	if err != nil {
		return err
	}
	bdb.dbConn = db
	return nil
}

/*
ConfirmSavantMaster makes sure that the table is present.  If not, create it.
*/
func (bdb *BaseballDB) ConfirmSavantMaster() error {
	_, err := bdb.dbConn.Exec(`create table if not exists mlb_savant (
		pitch_type text,
		game_date date,
		release_speed double precision,
		release_pos_x double precision,
		release_pos_z double precision,
		player_name text,
		batter integer,
		pitcher integer,
		events text,
		description text,
		spin_dir double precision,
		spin_rate_depricated double precision,
		break_angle_depricated double precision,
		break_length_depricated double precision,
		zone integer,
		des text,
		game_type text,
		stand text,
		p_throws text,
		home_team text,
		away_team text,
		type text,
		hit_location text,
		bb_type text,
		balls integer,
		strikes integer,
		game_year integer,
		pfx_x double precision,
		pfx_z double precision,
		plate_x double precision,
		plate_z double precision,
		on_3b integer,
		on_2b integer,
		on_1b integer,
		outs_when_up integer,
		inning integer,
		inning_topbot text,
		hc_x double precision,
		hc_y double precision,
		tfs_depricated text,
		tfs_zulu_depricated text,
		fielder_2 integer,
		umpire integer,
		sv_id text,
		vx0 double precision,
		vy0 double precision,
		vz0 double precision,
		ax double precision,
		ay double precision,
		az double precision,
		sz_top double precision,
		sz_bot double precision,
		hit_distance double precision,
		launch_speed double precision,
		launch_angle double precision,
		effective_speed double precision,
		release_spin double precision,
		release_extension double precision,
		game_pk integer,
		pitcher_id integer,
		catcher_id integer,
		firstbase_id integer,
		secondbase_id integer,
		thirdbase_id integer,
		shortstop_id integer,
		leftfield_id integer,
		centerfield_id integer,
		rightfield_id integer,
		release_pos_y double precision,
		estimated_ba_using_speedangle double precision,
		estimated_woba_using_speedangle double precision,
		woba_value double precision,
		woba_denom double precision,
		babip_value double precision,
		iso_value double precision,
		launch_speed_angle double precision,
		at_bat_number integer,
		pitch_number integer,
		pitch_name text,
		home_score integer,
		away_score integer,
		bat_score integer,
		fld_score integer,
		post_away_score integer,
		post_home_score integer,
		post_bat_score integer,
		if_fielding_alignment text,
		of_fielding_alignment text);
	`)
	if err != nil {
		return err
	}
	return nil
}

/*
LoadSavantCSV takes a CSV file and bulk loads it into the database
*/
func (bdb *BaseballDB) LoadSavantCSV(f string) error {
	fp, err := os.Open(f)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()
	r := csv.NewReader(fp)

	txn, err := bdb.dbConn.Begin()
	if err != nil {
		log.Fatal(err)
	}

	header, err := r.Read()
	if err != nil {
		log.Fatal(err)
	}
	dblRecord := make([]float64, len(header))
	intRecord := make([]int, len(header))

	stmt, err := txn.Prepare(pq.CopyIn("mlb_savant", "pitch_type",
		"game_date", "release_speed", "release_pos_x",
		"release_pos_z", "player_name",
		"batter", "pitcher", "events", "description",
		"spin_dir", "spin_rate_depricated",
		"break_angle_depricated", "break_length_depricated",
		"zone", "des", "game_type", "stand", "p_throws", "home_team",
		"away_team", "type", "hit_location", "bb_type",
		"balls", "strikes", "game_year", "pfx_x", "pfx_z",
		"plate_x", "plate_z", "on_3b", "on_2b", "on_1b",
		"outs_when_up", "inning", "inning_topbot", "hc_x",
		"hc_y", "tfs_depricated", "tfs_zulu_depricated",
		"fielder_2", "umpire", "sv_id", "vx0", "vy0", "vz0", "ax", "ay",
		"az", "sz_top", "sz_bot", "hit_distance", "launch_speed",
		"launch_angle", "effective_speed", "release_spin",
		"release_extension", "game_pk", "pitcher_id", "catcher_id",
		"firstbase_id", "secondbase_id", "thirdbase_id", "shortstop_id",
		"leftfield_id", "centerfield_id", "rightfield_id",
		"release_pos_y", "estimated_ba_using_speedangle",
		"estimated_woba_using_speedangle", "woba_value",
		"woba_denom", "babip_value", "iso_value", "launch_speed_angle",
		"at_bat_number", "pitch_number", "pitch_name", "home_score",
		"away_score", "bat_score", "fld_score", "post_away_score",
		"post_home_score", "post_bat_score", "if_fielding_alignment",
		"of_fielding_alignment"))
	if err != nil {
		return err
	}

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		for i, v := range record {
			if v == "null" {
				record[i] = ""
			}
			if strings.Contains(v, ",") {
				record[i] = strings.ReplaceAll(v, ",", "_")
			}
			switch i {
			case 2, 3, 4, 10, 11, 12, 13, 27, 28, 29, 30, 37, 38,
				44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55,
				56, 57, 68, 69, 70, 71, 72, 73, 74, 75:
				if len(record[i]) == 0 {
					dblRecord[i] = -88.0
				} else {
					dblRecord[i], err = strconv.ParseFloat(record[i], 64)
					if err != nil {
						log.Printf("%d -> %s\n", i, err)
						dblRecord[i] = -99.0
					}
				}
			case 6, 7, 14, 24, 25, 26, 31, 32, 33, 34, 35, 41, 42, 58,
				59, 60, 61, 62, 63, 64, 65, 66, 67, 76, 77, 79, 80, 81,
				82, 83, 84, 85:
				if len(record[i]) == 0 {
					intRecord[i] = -88
				} else {
					intRecord[i], err = strconv.Atoi(record[i])
					if err != nil {
						log.Printf("%d -> %s\n", i, err)
						intRecord[i] = -99
					}
				}
			}
		}
		_, err = stmt.Exec(record[0], record[1], dblRecord[2],
			dblRecord[3], dblRecord[4], record[5], intRecord[6], intRecord[7],
			record[8], record[9], dblRecord[10], dblRecord[11], dblRecord[12],
			dblRecord[13], intRecord[14], record[15], record[16], record[17],
			record[18], record[19], record[20], record[21], record[22],
			record[23], intRecord[24], intRecord[25], intRecord[26], dblRecord[27],
			dblRecord[28], dblRecord[29], dblRecord[30], intRecord[31], intRecord[32],
			intRecord[33], intRecord[34], intRecord[35], record[36], dblRecord[37],
			dblRecord[38], record[39], record[40], intRecord[41], intRecord[42],
			record[43], dblRecord[44], dblRecord[45], dblRecord[46], dblRecord[47],
			dblRecord[48], dblRecord[49], dblRecord[50], dblRecord[51], dblRecord[52],
			dblRecord[53], dblRecord[54], dblRecord[55], dblRecord[56], dblRecord[57],
			intRecord[58], intRecord[59], intRecord[60], intRecord[61], intRecord[62],
			intRecord[63], intRecord[64], intRecord[65], intRecord[66], intRecord[67],
			dblRecord[68], dblRecord[69], dblRecord[70], dblRecord[71], dblRecord[72],
			dblRecord[73], dblRecord[74], dblRecord[75], intRecord[76], intRecord[77],
			record[78], intRecord[79], intRecord[80], intRecord[81], intRecord[82],
			intRecord[83], intRecord[84], intRecord[85], record[86], record[87])
		if err != nil {
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	err = stmt.Close()
	if err != nil {
		return err
	}

	err = txn.Commit()
	if err != nil {
		return err
	}
	return nil
}

/*
Close closes the connection to the database
*/
func (bdb *BaseballDB) Close() error {
	return bdb.dbConn.Close()
}

/*
DropSavantTable gets rid of the table if it exists
*/
func (bdb *BaseballDB) DropSavantTable() error {
	_, err := bdb.dbConn.Exec("drop table if exists mlb_savant;")
	return err
}

/*
ConfirmGamedayMaster makes sure that the table is present.  If not, create it.
*/
func (bdb *BaseballDB) ConfirmGamedayMaster() error {
	_, err := bdb.dbConn.Exec(`create table if not exists mlb_gameday (
		game_date date,
		away_team text,
		home_team text,
		game_number integer,
		inning integer,
		at_bat_number integer,
		at_bat_start_tfs integer,
		at_bat_start_tfs_zulu timestamp with time zone,
		at_bat_end_tfs_zulu timestamp with time zone,
		pitch_number integer,
		sv_id text,
		pitch_tfs integer,
		pitch_tfs_zulu timestamp with time zone);
		`)
	return err
}

/*
DropGamedayTable gets rid of the table if it exists
*/
func (bdb *BaseballDB) DropGamedayTable() error {
	_, err := bdb.dbConn.Exec("drop table if exists mlb_gameday;")
	return err
}

/*
<pitch
break_angle="38.4"
break_length="4.8"
break_y="24.0"
cc=""
code="B"
des="Ball"
des_es="In play, out(s)"
end_speed="84.1"
event_num="2"
id="2"
mt=""
nasty=""
pitch_type="FF"
play_guid="28c604b5-9faf-43db-b371-7535480c9a4e"
spin_dir="placeholder"
spin_rate="placeholder"
start_speed="92.5"
*/

/*
PitchXML represents a pitch
*/
type PitchXML struct {
	SVID           string `xml:"sv_id,attr"`
	TFS            string `xml:"tfs,attr"`
	TFSZulu        string `xml:"tfs_zulu,attr"`
	X              string `xml:"x,attr"`
	Y              string `xml:"y,attr"`
	X0             string `xml:"x0,attr"`
	Y0             string `xml:"y0,attr"`
	Z0             string `xml:"z0,attr"`
	VX0            string `xml:"vx0,attr"`
	VY0            string `xml:"vy0,attr"`
	VZ0            string `xml:"vz0,attr"`
	AX0            string `xml:"ax0,attr"`
	AY0            string `xml:"ay0,attr"`
	AZ0            string `xml:"az0,attr"`
	PX             string `xml:"px,attr"`
	PZ             string `xml:"pz,attr"`
	PFXX           string `xml:"pfx_x,attr"`
	PFXZ           string `xml:"pfx_z,attr"`
	Type           string `xml:"type,attr"`
	TypeConfidence string `xml:"type_confidence"`
	Zone           string `xml:"zone,attr"`
	SZTop          string `xml:"sz_top,attr"`
	SZBot          string `xml:"sz_bot,attr"`
}

/*
AtBatXML represents an at bat
*/
type AtBatXML struct {
	PlayNumber       int        `xml:"num,attr"`
	PlayGUID         string     `xml:"play_guid,attr"`
	AwayTeamRuns     int        `xml:"away_team_runs,attr"`
	HomeTeamRuns     int        `xml:"home_team_runs,attr"`
	Balls            int        `xml:"b,attr"`
	Strikes          int        `xml:"s,attr"`
	Outs             int        `xml:"o,attr"`
	BatterID         int        `xml:"batter,attr"`
	BatterHeight     string     `xml:"b_height,attr"`
	BatterSide       string     `xml:"stand,attr"`
	PitcherID        int        `xml:"pitcher,attr"`
	PitcherSide      string     `xml:"p_throws,attr"`
	EnglishDesc      string     `xml:"des,attr"`
	SpanishDesc      string     `xml:"des_es,attr"`
	EnglishEventDesc string     `xml:"event,attr"`
	SpanishEventDesc string     `xml:"event_es,attr"`
	EventNumber      int        `xml:"event_num,attr"`
	StartTFSZulu     string     `xml:"start_tfs_zulu,attr"`
	EndTFSZulu       string     `xml:"end_tfs_zulu,attr"`
	StartTFS         string     `xml:"start_tfs,attr"`
	Pitches          []PitchXML `xml:"pitch"`
}

/*
HalfInningXML represents a half inning
*/
type HalfInningXML struct {
	AtBats []AtBatXML `xml:"atbat"`
}

/*
InningXML represents an inning
*/
type InningXML struct {
	AwayTeam string        `xml:"away_team,attr"`
	HomeTeam string        `xml:"home_team,attr"`
	Next     string        `xml:"next,attr"`
	Num      int           `xml:"num,attr"`
	Top      HalfInningXML `xml:"top"`
	Bottom   HalfInningXML `xml:"bottom"`
}

/*
GameXML is the highest level in the xml file
*/
type GameXML struct {
	AtBat   int         `xml:"atBat,attr"`
	Deck    int         `xml:"deck,attr"`
	Hole    int         `xml:"hole,attr"`
	Ind     string      `xml:"ind,attr"`
	Innings []InningXML `xml:"inning"`
}

/*
LoadGamedayXML takes an XML file and bulk loads it into the database
*/
func (bdb *BaseballDB) LoadGamedayXML(f string) error {
	fp, err := os.Open(f)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()

	fileComponents := strings.Split(f, "_")
	gY := fileComponents[1]
	gM := fileComponents[2]
	gD := fileComponents[3]
	gN := fileComponents[6]

	var g GameXML
	decoder := xml.NewDecoder(fp)
	err = decoder.Decode(&g)
	if err != nil {
		return err
	}

	txn, err := bdb.dbConn.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := txn.Prepare(pq.CopyIn("mlb_gameday",
		"game_date", "away_team", "home_team", "game_number",
		"inning", "at_bat_number", "at_bat_start_tfs",
		"at_bat_start_tfs_zulu", "at_bat_end_tfs_zulu",
		"pitch_number", "sv_id", "pitch_tfs",
		"pitch_tfs_zulu"))
	if err != nil {
		return err
	}

	atBatNum := 0
	for _, inning := range g.Innings {
		for _, atbat := range inning.Top.AtBats {
			atBatNum++
			for pitchnum, pitch := range atbat.Pitches {
				gameDate := fmt.Sprintf("%s-%s-%s", gY, gM, gD)
				_, err = stmt.Exec(gameDate, strings.ToUpper(inning.AwayTeam),
					strings.ToUpper(inning.HomeTeam), gN, inning.Num,
					atBatNum, atbat.StartTFS, atbat.StartTFSZulu, atbat.EndTFSZulu,
					pitchnum+1, pitch.SVID, pitch.TFS, pitch.TFSZulu)
				if err != nil {
					return err
				}
			}
		}
		for _, atbat := range inning.Bottom.AtBats {
			atBatNum++
			for pitchnum, pitch := range atbat.Pitches {
				gameDate := fmt.Sprintf("%s-%s-%s", gY, gM, gD)
				_, err = stmt.Exec(gameDate, strings.ToUpper(inning.AwayTeam),
					strings.ToUpper(inning.HomeTeam), gN, inning.Num,
					atBatNum, atbat.StartTFS, atbat.StartTFSZulu, atbat.EndTFSZulu,
					pitchnum+1, pitch.SVID, pitch.TFS, pitch.TFSZulu)
				if err != nil {
					return err
				}
			}
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	err = stmt.Close()
	if err != nil {
		return err
	}

	err = txn.Commit()
	if err != nil {
		return err
	}

	return nil
}
