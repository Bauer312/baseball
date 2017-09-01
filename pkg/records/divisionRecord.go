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
DivisionRecord is the specific data record for each division in a league
*/
type DivisionRecord struct {
	RecordName    string
	EffectiveDate time.Time
	Name          string
	Code          string
}

/*
ScreenOutput displays the record on the screen
*/
func (dR *DivisionRecord) ScreenOutput() {
	fmt.Println(dR)
}

/*
FileOutput displays the record on the screen
*/
func (dR *DivisionRecord) FileOutput(filePtr *os.File) {
	fmt.Fprintf(filePtr, "%s|%s|%s\n",
		dR.EffectiveDate.Format(time.UnixDate),
		dR.Name,
		dR.Code,
	)
}

/*
CreateTable will create the requisite database table
*/
func (dR *DivisionRecord) CreateTable(db *sql.DB) {
	statement := `CREATE TABLE IF NOT EXISTS DivisionRecord (
		effectiveDate 	timestamp with time zone,
		name 			varchar(128),
		code 			varchar(16),
		PRIMARY KEY (name, code)
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
func (dR *DivisionRecord) UpdateRecord(db *sql.DB) {
	/*
		1.  If this is a unique record, insert it.
		2.  If this is a duplicate record and the effective date is earlier,
				update the existing record.
	*/
	statement := `INSERT INTO DivisionRecord VALUES ($1,$2,$3);`
	_, err := db.Exec(statement, dR.EffectiveDate, dR.Name, dR.Code)
	if err != nil {
		if pqerr, ok := err.(*pq.Error); ok {
			if pqerr.Code.Name() == "unique_violation" {
				var existingEffectiveDate time.Time
				statement = `SELECT effectiveDate FROM DivisionRecord WHERE
				name=$1 AND code=$2;`
				err = db.QueryRow(statement, dR.Name, dR.Code).Scan(&existingEffectiveDate)
				if err != nil {
					if pqerr, ok := err.(*pq.Error); ok {
						fmt.Println("pq error:", pqerr.Code.Name())
					} else {
						fmt.Println(err)
					}
				}
				if existingEffectiveDate.Sub(dR.EffectiveDate) > 0 {
					//The new date is before the existing date, so update the record in the database
					//fmt.Printf("Existing: %v New: %v --> Updating record in DB\n", existingEffectiveDate, dR.EffectiveDate)
					statement = `UPDATE DivisionRecord SET effectiveDate=$1 WHERE
					name=$2 AND code=$3;`
					_, err := db.Exec(statement, dR.EffectiveDate, dR.Name, dR.Code)
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