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

package records

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	pq "github.com/lib/pq"
)

/*
GameRecord is the specific data record for a game
*/
type GameRecord struct {
	RecordName       string
	EffectiveDate    time.Time
	ID               int64
	ResumeDate       string
	OriginalDate     string
	GameType         string
	Tiebreaker       string
	GameDay          string
	DoubleHeader     string
	GameNumber       int
	TBDFlag          string
	Interleague      string
	ScheduledInnings int
	Description      string
	VenueID          int64
	AwayTeamID       int64
	HomeTeamID       int64
}

/*
ScreenOutput displays the record on the screen
*/
func (gR *GameRecord) ScreenOutput() {
	fmt.Println(gR)
}

/*
FileOutput displays the record on the screen
*/
func (gR *GameRecord) FileOutput(filePtr *os.File) {
	fmt.Fprintf(filePtr, "%s|%d|%s|%s|%s|%s|%s|%s|%d|%s|%s|%d|%s|%d|%d|%d\n",
		gR.EffectiveDate.Format(time.UnixDate),
		gR.ID,
		gR.ResumeDate,
		gR.OriginalDate,
		gR.GameType,
		gR.Tiebreaker,
		gR.GameDay,
		gR.DoubleHeader,
		gR.GameNumber,
		gR.TBDFlag,
		gR.Interleague,
		gR.ScheduledInnings,
		gR.Description,
		gR.VenueID,
		gR.AwayTeamID,
		gR.HomeTeamID,
	)
}

/*
CreateTable will create the requisite database table
*/
func (gR *GameRecord) CreateTable(db *sql.DB) {
	statement := `CREATE TABLE IF NOT EXISTS GameRecord (
		effectiveDate 		timestamp with time zone,
		id 					bigint,
		resumedate			varchar(64),
		originaldate		varchar(64),
		gametype			varchar(8),
		tiebreaker			varchar(8),
		gameday				varchar(8),
		doubleheader		varchar(8),
		gamenumber			int,
		tbdflag				varchar(8),
		interleague			varchar(8),
		scheduledinnings	int,
		description			varchar(256),
		venueid				bigint,
		awayteamid			bigint,
		hometeamid			bigint,
		PRIMARY KEY (id)
	)`

	_, err := db.Exec(statement)
	if err != nil {
		fmt.Println(err)
	}
}

/*
UpdateRecord is the way data gets into the database.  It does not act like
	the UPSERT command because the effective date field will be different
	for each record.  Each table in the database will have different rules
	for how to deal with data records
*/
func (gR *GameRecord) UpdateRecord(db *sql.DB) {
	/*
		1.  If this is a unique record, insert it.
		2.  If this is a duplicate record and the effective date is later,
				update the existing record.
	*/
	statement := `INSERT INTO GameRecord VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16);`
	_, err := db.Exec(statement, gR.EffectiveDate.UTC(), gR.ID, gR.ResumeDate, gR.OriginalDate, gR.GameType,
		gR.Tiebreaker, gR.GameDay, gR.DoubleHeader, gR.GameNumber, gR.TBDFlag, gR.Interleague,
		gR.ScheduledInnings, gR.Description, gR.VenueID, gR.AwayTeamID, gR.HomeTeamID)
	if err != nil {
		if pqerr, ok := err.(*pq.Error); ok {
			if pqerr.Code.Name() == "unique_violation" {
				var existingEffectiveDate time.Time
				statement = `SELECT effectiveDate FROM GameRecord WHERE id=$1;`
				err = db.QueryRow(statement, gR.ID).Scan(&existingEffectiveDate)
				if err != nil {
					if pqerr, ok := err.(*pq.Error); ok {
						fmt.Println("pq error:", pqerr.Code.Name())
					} else {
						fmt.Println(err)
					}
				}
				if gR.EffectiveDate.UTC().Sub(existingEffectiveDate) > 0 {
					//The new date is after the existing date, so update the record in the database
					//fmt.Printf("Existing: %v New: %v --> Replacing record in DB\n", existingEffectiveDate, gR.EffectiveDate)
					statement = `UPDATE GameRecord SET effectiveDate=$1, resumedate=$2, originaldate=$3,
					gametype=$4, tiebreaker=$5, gameday=$6, doubleheader=$7, gamenumber=$8, tbdflag=$9,
					interleague=$10, scheduledinnings=$11, description=$12, venueid=$13, awayteamid=$14,
					hometeamid=$15 WHERE id=$16;`
					_, err := db.Exec(statement, gR.EffectiveDate.UTC(), gR.ResumeDate, gR.OriginalDate, gR.GameType,
						gR.Tiebreaker, gR.GameDay, gR.DoubleHeader, gR.GameNumber, gR.TBDFlag, gR.Interleague,
						gR.ScheduledInnings, gR.Description, gR.VenueID, gR.AwayTeamID, gR.HomeTeamID, gR.ID)
					if err != nil {
						if pqerr, ok := err.(*pq.Error); ok {
							fmt.Println("pq error:", pqerr.Code.Name())
						} else {
							fmt.Println(err)
						}
					}
				}
			} else {
				fmt.Println("pq error:", pqerr.Code.Name())
			}
		} else {
			fmt.Println(err)
		}
	}
}
