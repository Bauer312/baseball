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
LeagueRecord is the specific data record for each league
*/
type LeagueRecord struct {
	RecordName    string
	EffectiveDate time.Time
	ID            int64
	Name          string
	SportCode     string
}

/*
ScreenOutput displays the record on the screen
*/
func (lR *LeagueRecord) ScreenOutput() {
	fmt.Println(lR)
}

/*
FileOutput displays the record on the screen
*/
func (lR *LeagueRecord) FileOutput(filePtr *os.File) {
	fmt.Fprintf(filePtr, "%s|%d|%s|%s\n",
		lR.EffectiveDate.Format(time.UnixDate),
		lR.ID,
		lR.Name,
		lR.SportCode,
	)
}

/*
CreateTable will create the requisite database table
*/
func (lR *LeagueRecord) CreateTable(db *sql.DB) {
	statement := `CREATE TABLE IF NOT EXISTS LeagueRecord (
		effectiveDate 	timestamp with time zone,
		id 				bigint,
		name 			varchar(128),
		sportCode 		varchar(16),
		PRIMARY KEY (id, name, sportCode)
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
func (lR *LeagueRecord) UpdateRecord(db *sql.DB) {
	/*
		1.  If this is a unique record, insert it.
		2.  If this is a duplicate record and the effective date is earlier,
				update the existing record.
	*/
	statement := `INSERT INTO LeagueRecord VALUES ($1,$2,$3,$4);`
	_, err := db.Exec(statement, lR.EffectiveDate, lR.ID, lR.Name, lR.SportCode)
	if err != nil {
		if pqerr, ok := err.(*pq.Error); ok {
			if pqerr.Code.Name() == "unique_violation" {
				var existingEffectiveDate time.Time
				statement = `SELECT effectiveDate FROM LeagueRecord WHERE
				id=$1 AND name=$2 AND sportCode=$3;`
				err = db.QueryRow(statement, lR.ID, lR.Name, lR.SportCode).Scan(&existingEffectiveDate)
				if err != nil {
					if pqerr, ok := err.(*pq.Error); ok {
						fmt.Println("pq error:", pqerr.Code.Name())
					} else {
						fmt.Println(err)
					}
				}
				if existingEffectiveDate.Sub(lR.EffectiveDate) > 0 {
					//The new date is before the existing date, so update the record in the database
					//fmt.Printf("Existing: %v New: %v --> Updating record in DB\n", existingEffectiveDate, lR.EffectiveDate)
					statement = `UPDATE LeagueRecord SET effectiveDate=$1 WHERE
					id=$2 AND name=$3 AND sportCode=$4;`
					_, err := db.Exec(statement, lR.EffectiveDate, lR.ID, lR.Name, lR.SportCode)
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
