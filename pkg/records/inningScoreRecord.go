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
InningScoreRecord is the specific data about an inning in a game
*/
type InningScoreRecord struct {
	RecordName    string
	EffectiveDate time.Time
	GameID        int64
	Inning        int
	AwayTeamRuns  int
	HomeTeamRuns  int
}

/*
ScreenOutput displays the record on the screen
*/
func (isR *InningScoreRecord) ScreenOutput() {
	fmt.Println(isR)
}

/*
FileOutput displays the record on the screen
*/
func (isR *InningScoreRecord) FileOutput(filePtr *os.File) {
	fmt.Fprintf(filePtr, "%s|%d|%d|%d|%d\n",
		isR.EffectiveDate.Format(time.UnixDate),
		isR.GameID,
		isR.Inning,
		isR.AwayTeamRuns,
		isR.HomeTeamRuns,
	)
}

/*
CreateTable will create the requisite database table
*/
func (isR *InningScoreRecord) CreateTable(db *sql.DB) {
	statement := `CREATE TABLE IF NOT EXISTS InningScoreRecord (
		effectiveDate 	timestamp with time zone,
		gameid			bigint,
		inning			int,
		awayteamruns	int,
		hometeamruns	int,
		PRIMARY KEY (gameid, inning)
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
func (isR *InningScoreRecord) UpdateRecord(db *sql.DB) {
	/*
		1.  If this is a unique record, insert it.
		2.  If this is a duplicate record and the effective date is later,
				replace the existing record.
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
	statement = `INSERT INTO InningScoreRecord VALUES ($1,$2,$3,$4,$5);`
	_, err = db.Exec(statement, isR.EffectiveDate.UTC(), isR.GameID, isR.Inning, isR.AwayTeamRuns, isR.HomeTeamRuns)
	if err != nil {
		if pqerr, ok := err.(*pq.Error); ok {
			if pqerr.Code.Name() == "unique_violation" {
				var existingEffectiveDate time.Time
				statement = `SELECT effectiveDate FROM InningScoreRecord WHERE
				gameid=$1 AND inning=$2;`
				err = db.QueryRow(statement, isR.GameID, isR.Inning).Scan(&existingEffectiveDate)
				if err != nil {
					if pqerr, ok := err.(*pq.Error); ok {
						fmt.Println("pq error:", pqerr.Code.Name())
					} else {
						fmt.Println(err)
					}
				}
				if isR.EffectiveDate.UTC().Sub(existingEffectiveDate) > 0 {
					//The new date is after the existing date, so replace the record in the database
					//fmt.Printf("Existing: %v New: %v --> Replacing record in DB\n", existingEffectiveDate, isR.EffectiveDate)
					statement = `UPDATE InningScoreRecord SET effectiveDate=$1, awayteamruns=$2, hometeamruns=$3 WHERE
					gameid=$4 AND inning=$5;`
					_, err := db.Exec(statement, isR.EffectiveDate.UTC(), isR.AwayTeamRuns, isR.HomeTeamRuns, isR.GameID, isR.Inning)
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
