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

	//Make sure we can use Postgres
	pq "github.com/lib/pq"
)

/*
VenueRecord is the specific data record for each venue
*/
type VenueRecord struct {
	RecordName    string
	EffectiveDate time.Time
	ID            int64
	Name          string
	Location      string
	Channel       string
}

/*
ScreenOutput displays the record on the screen
*/
func (vR *VenueRecord) ScreenOutput() {
	fmt.Println(vR)
}

/*
FileOutput displays the record on the screen
*/
func (vR *VenueRecord) FileOutput(filePtr *os.File) {
	fmt.Fprintf(filePtr, "%s|%d|%s|%s|%s\n",
		vR.EffectiveDate.Format(time.UnixDate),
		vR.ID,
		vR.Name,
		vR.Location,
		vR.Channel,
	)
}

/*
CreateTable will create the requisite database table
*/
func (vR *VenueRecord) CreateTable(db *sql.DB) {
	statement := `CREATE TABLE IF NOT EXISTS VenueRecord (
		effectiveDate 	timestamp with time zone,
		id 				bigint,
		name 			varchar(128),
		location 		varchar(128),
		channel 		varchar(16),
		PRIMARY KEY (id, name, location, channel)
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
func (vR *VenueRecord) UpdateRecord(db *sql.DB) {
	/*
		1.  If this is a unique record, insert it.
		2.  If this is a duplicate record and the effective date is earlier,
				update the existing record.
	*/
	statement := `INSERT INTO VenueRecord VALUES ($1,$2,$3,$4,$5);`
	_, err := db.Exec(statement, vR.EffectiveDate, vR.ID, vR.Name, vR.Location, vR.Channel)
	if err != nil {
		if pqerr, ok := err.(*pq.Error); ok {
			if pqerr.Code.Name() == "unique_violation" {
				var existingEffectiveDate time.Time
				statement = `SELECT effectiveDate FROM VenueRecord WHERE
				id=$1 AND name=$2 AND location=$3 AND channel=$4;`
				err = db.QueryRow(statement, vR.ID, vR.Name, vR.Location, vR.Channel).Scan(&existingEffectiveDate)
				if err != nil {
					if pqerr, ok := err.(*pq.Error); ok {
						fmt.Println("pq error:", pqerr.Code.Name())
					} else {
						fmt.Println(err)
					}
				}
				if existingEffectiveDate.Sub(vR.EffectiveDate) > 0 {
					//The new date is before the existing date, so update the record in the database
					//fmt.Printf("Existing: %v New: %v --> Updating record in DB\n", existingEffectiveDate, vR.EffectiveDate)
					statement = `UPDATE VenueRecord SET effectiveDate=$1 WHERE
					id=$2 AND name=$3 AND location=$4 AND channel=$5;`
					_, err := db.Exec(statement, vR.EffectiveDate, vR.ID, vR.Name, vR.Location, vR.Channel)
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
