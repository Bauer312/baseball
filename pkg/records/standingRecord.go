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
StandingRecord is the specific data record for a team standing in
	its division in its league
*/
type StandingRecord struct {
	RecordName        string
	EffectiveDate     time.Time
	TeamID            int64
	Wins              int
	Losses            int
	GamesBack         string
	WildcardGamesBack string
}

/*
ScreenOutput displays the record on the screen
*/
func (sR *StandingRecord) ScreenOutput() {
	fmt.Println(sR)
}

/*
FileOutput displays the record on the screen
*/
func (sR *StandingRecord) FileOutput(filePtr *os.File) {
	fmt.Fprintf(filePtr, "%s|%d|%d|%d|%s|%s\n",
		sR.EffectiveDate.Format(time.UnixDate),
		sR.TeamID,
		sR.Wins,
		sR.Losses,
		sR.GamesBack,
		sR.WildcardGamesBack,
	)
}

/*
CreateTable will create the requisite database table
*/
func (sR *StandingRecord) CreateTable(db *sql.DB) {
	statement := `CREATE TABLE IF NOT EXISTS StandingRecord (
		effectiveDate 		timestamp with time zone,
		teamid 				bigint,
		wins				int,
		losses				int,
		gamesplayed			int,
		gamesback			varchar(8),
		wildcardgamesback	varchar(8),
		PRIMARY KEY (effectiveDate, teamid)
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
func (sR *StandingRecord) UpdateRecord(db *sql.DB) {
	/*
		1.  If this is a unique record, insert it.
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
	statement = `INSERT INTO StandingRecord VALUES ($1,$2,$3,$4,$5,$6,$7);`
	_, err = db.Exec(statement, sR.EffectiveDate.UTC(), sR.TeamID, sR.Wins, sR.Losses, sR.Wins+sR.Losses, sR.GamesBack, sR.WildcardGamesBack)
	if err != nil {
		if pqerr, ok := err.(*pq.Error); ok {
			if pqerr.Code.Name() == "unique_violation" {
				var existingGamesPlayed int
				statement = `SELECT gamesplayed FROM StandingRecord WHERE effectivedate = $1 AND teamid = $2;`
				err = db.QueryRow(statement, sR.EffectiveDate.UTC(), sR.TeamID).Scan(&existingGamesPlayed)
				if err != nil {
					if pqerr, ok := err.(*pq.Error); ok {
						fmt.Println("pq error:", pqerr.Code.Name())
					} else {
						fmt.Println(err)
					}
				}
				if sR.Wins+sR.Losses > existingGamesPlayed {
					statement = `UPDATE StandingRecord SET
					wins = $1, losses = $2, gamesplayed = $3, gamesback = $4, wildcardgamesback = $5
					WHERE effectivedate = $6 AND teamid = $7;`
					_, err = db.Exec(statement, sR.Wins, sR.Losses, sR.GamesBack, sR.WildcardGamesBack, sR.Wins+sR.Losses, sR.EffectiveDate.UTC(), sR.TeamID)
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
