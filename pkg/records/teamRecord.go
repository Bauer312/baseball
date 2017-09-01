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
TeamRecord is the specifc data record for each team in a league
*/
type TeamRecord struct {
	RecordName    string
	EffectiveDate time.Time
	ID            int64
	Name          string
	Code          string
	City          string
	LeagueID      int64
	Division      string
}

/*
ScreenOutput displays the record on the screen
*/
func (tR *TeamRecord) ScreenOutput() {
	fmt.Println(tR)
}

/*
FileOutput displays the record on the screen
*/
func (tR *TeamRecord) FileOutput(filePtr *os.File) {
	fmt.Fprintf(filePtr, "%s|%d|%s|%s|%s|%d|%s\n",
		tR.EffectiveDate.Format(time.UnixDate),
		tR.ID,
		tR.Name,
		tR.Code,
		tR.City,
		tR.LeagueID,
		tR.Division,
	)
}

/*
CreateTable will create the requisite database table
*/
func (tR *TeamRecord) CreateTable(db *sql.DB) {
	statement := `CREATE TABLE IF NOT EXISTS TeamRecord (
		effectiveDate	timestamp with time zone,
		id 				bigint,
		name 			varchar(128),
		code			varchar(16),
		city	 		varchar(128),
		leagueid		bigint,
		division		varchar(32),
		PRIMARY KEY (id, name, code, city, leagueid, division)
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
func (tR *TeamRecord) UpdateRecord(db *sql.DB) {
	/*
		1.  If this is a unique record, insert it.
		2.  If this is a duplicate record and the effective date is earlier,
				update the existing record.
	*/
	statement := `INSERT INTO TeamRecord VALUES ($1, $2, $3, $4, $5, $6, $7);`
	_, err := db.Exec(statement, tR.EffectiveDate, tR.ID, tR.Name, tR.Code, tR.City, tR.LeagueID, tR.Division)
	if err != nil {
		if pqerr, ok := err.(*pq.Error); ok {
			if pqerr.Code.Name() == "unique_violation" {
				var existingEffectiveDate time.Time
				statement = `SELECT effectiveDate FROM TeamRecord WHERE
				id=$1 AND name=$2 AND code=$3 AND city=$4 AND leagueid=$5 AND division=$6;`
				err = db.QueryRow(statement, tR.ID, tR.Name, tR.Code, tR.City, tR.LeagueID, tR.Division).Scan(&existingEffectiveDate)
				if err != nil {
					if pqerr, ok := err.(*pq.Error); ok {
						fmt.Println("pq error:", pqerr.Code.Name())
					} else {
						fmt.Println(err)
					}
				}
				if existingEffectiveDate.Sub(tR.EffectiveDate) > 0 {
					//The new date is before the existing date, so update the record in the database
					//fmt.Printf("Existing: %v New: %v --> Updating record in DB\n", existingEffectiveDate, tR.EffectiveDate)
					statement = `UPDATE TeamRecord SET effectiveDate=$1 WHERE
					id=$2 AND name=$3 AND code=$4 AND city=$5 AND leagueid=$6 AND division=$7;`
					_, err := db.Exec(statement, tR.EffectiveDate, tR.ID, tR.Name, tR.Code, tR.City, tR.LeagueID, tR.Division)
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
