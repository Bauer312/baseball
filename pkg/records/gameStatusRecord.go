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
GameStatusRecord is the specific data about the game
*/
type GameStatusRecord struct {
	RecordName     string
	EffectiveDate  time.Time
	ID             int64
	Status         string
	Ind            string
	Reason         string
	CurrentInning  int
	TopOfInning    bool
	Balls          int
	Strikes        int
	Outs           int
	InningState    string
	Note           string
	PerfectGame    bool
	NoHitter       bool
	AwayTeamRuns   int
	HomeTeamRuns   int
	AwayTeamHits   int
	HomeTeamHits   int
	AwayTeamErrors int
	HomeTeamErrors int
	AwayTeamHR     int
	HomeTeamHR     int
	AwayTeamSB     int
	HomeTeamSB     int
	AwayTeamSO     int
	HomeTeamSO     int
	Innings        []InningScoreRecord
}

/*
ScreenOutput displays the record on the screen
*/
func (gsR *GameStatusRecord) ScreenOutput() {
	fmt.Println(gsR)
}

/*
FileOutput displays the record on the screen
*/
func (gsR *GameStatusRecord) FileOutput(filePtr *os.File) {
	fmt.Fprintf(filePtr, "%s|%d|%s|%s|%s|%d|%t|%d|%d|%d|%s|%s|%t|%t|%d|%d|%d|%d|%d|%d|%d|%d|%d|%d|%d|%d\n",
		gsR.EffectiveDate.Format(time.UnixDate),
		gsR.ID,
		gsR.Status,
		gsR.Ind,
		gsR.Reason,
		gsR.CurrentInning,
		gsR.TopOfInning,
		gsR.Balls,
		gsR.Strikes,
		gsR.Outs,
		gsR.InningState,
		gsR.Note,
		gsR.PerfectGame,
		gsR.NoHitter,
		gsR.AwayTeamRuns,
		gsR.HomeTeamRuns,
		gsR.AwayTeamHits,
		gsR.HomeTeamHits,
		gsR.AwayTeamErrors,
		gsR.HomeTeamErrors,
		gsR.AwayTeamHR,
		gsR.HomeTeamHR,
		gsR.AwayTeamSB,
		gsR.HomeTeamSB,
		gsR.AwayTeamSO,
		gsR.HomeTeamSO,
	)
}

/*
CreateTable will create the requisite database table
*/
func (gsR *GameStatusRecord) CreateTable(db *sql.DB) {
	statement := `CREATE TABLE IF NOT EXISTS GameStatusRecord (
		effectiveDate 	timestamp with time zone,
		id 				bigint,
		status			varchar(128),
		ind				varchar(8),
		reason			varchar(128),
		currentInning	int,
		topOfInning		boolean,
		balls			int,
		strikes			int,
		outs			int,
		inningState		varchar(8),
		note			varchar(128),
		perfectGame		boolean,
		noHitter		boolean,
		awayTeamRuns	int,
		homeTeamRuns	int,
		awayTeamHits	int,
		homeTeamHits	int,
		awayTeamErrors	int,
		homeTeamErrors	int,
		awayTeamHR		int,
		homeTeamHR		int,
		awayTeamSB		int,
		homeTeamSB		int,
		awayTeamSO		int,
		homeTeamSO		int,
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
func (gsR *GameStatusRecord) UpdateRecord(db *sql.DB) {
	/*
		1.  If this is a unique record, insert it.
		2.  If this is a duplicate record and the effective date is later,
				update the existing record.
	*/
	statement := `SET timezone='UTC';`
	_, err := db.Exec(statement)
	if err != nil {
		if pqerr, ok := err.(*pq.Error); ok {
			fmt.Println("pq error:", pqerr.Code.Name())
		} else {
			fmt.Println(err)
		}
	}
	statement = `INSERT INTO GameStatusRecord VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,
	$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26);`
	_, err = db.Exec(statement, gsR.EffectiveDate.UTC(), gsR.ID, gsR.Status, gsR.Ind, gsR.Reason,
		gsR.CurrentInning, gsR.TopOfInning, gsR.Balls, gsR.Strikes, gsR.Outs, gsR.InningState,
		gsR.Note, gsR.PerfectGame, gsR.NoHitter, gsR.AwayTeamRuns, gsR.HomeTeamRuns,
		gsR.AwayTeamHits, gsR.HomeTeamHits, gsR.AwayTeamErrors, gsR.HomeTeamErrors,
		gsR.AwayTeamHR, gsR.HomeTeamHR, gsR.AwayTeamSB, gsR.HomeTeamSB,
		gsR.AwayTeamSO, gsR.HomeTeamSO)
	if err != nil {
		if pqerr, ok := err.(*pq.Error); ok {
			if pqerr.Code.Name() == "unique_violation" {
				var existingEffectiveDate time.Time
				statement = `SELECT effectiveDate FROM GameStatusRecord WHERE id=$1;`
				err = db.QueryRow(statement, gsR.ID).Scan(&existingEffectiveDate)
				if err != nil {
					if pqerr, ok := err.(*pq.Error); ok {
						fmt.Println("pq error:", pqerr.Code.Name())
					} else {
						fmt.Println(err)
					}
				}
				if gsR.EffectiveDate.UTC().Sub(existingEffectiveDate) > 0 {
					//The new date is after the existing date, so update the record in the database
					//fmt.Printf("Existing: %v New: %v --> Replacing record in DB\n", existingEffectiveDate, gsR.EffectiveDate)
					statement = `UPDATE GameStatusRecord SET effectiveDate=$1, status=$2, ind=$3, reason=$4,
					currentInning=$5, topOfInning=$6, balls=$7, strikes=$8, outs=$9, inningState=$10, note=$11,
					perfectGame=$12, noHitter=$13, awayTeamRuns=$14, homeTeamRuns=$15, awayTeamHits=$16,
					homeTeamHits=$17, awayTeamErrors=$18, homeTeamErrors=$19, awayTeamHR=$20, homeTeamHR=$21,
					awayTeamSB=$22, homeTeamSB=$23, awayTeamSO=$24, homeTeamSO=$25 WHERE id=$26`
					_, err := db.Exec(statement, gsR.EffectiveDate.UTC(), gsR.Status, gsR.Ind, gsR.Reason, gsR.CurrentInning,
						gsR.TopOfInning, gsR.Balls, gsR.Strikes, gsR.Outs, gsR.InningState, gsR.Note, gsR.PerfectGame,
						gsR.NoHitter, gsR.AwayTeamRuns, gsR.HomeTeamRuns, gsR.AwayTeamHits, gsR.HomeTeamHits,
						gsR.AwayTeamErrors, gsR.HomeTeamErrors, gsR.AwayTeamHR, gsR.HomeTeamHR, gsR.AwayTeamSB,
						gsR.HomeTeamSB, gsR.AwayTeamSO, gsR.HomeTeamSO, gsR.ID)
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
