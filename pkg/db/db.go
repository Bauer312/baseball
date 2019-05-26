package db

import (
	"database/sql"
	"encoding/csv"
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
