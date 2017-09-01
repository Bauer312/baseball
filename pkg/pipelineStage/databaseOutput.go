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

package pipelineStage

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	records "github.com/bauer312/baseball/pkg/records"
	//Make sure we can use Postgres
	_ "github.com/lib/pq"
)

/*
DatabaseOutput contains the elements of a pipeline stage that will accept
	strings of data and print them to the screen
*/
type DatabaseOutput struct {
	DataInput []chan string
	wg        sync.WaitGroup
	db        *sql.DB
	tables    map[string]bool
}

/*
Init the pipeline stage,
*/
func (dbO *DatabaseOutput) Init() error {
	numChannels := len(dbO.DataInput)
	dbO.wg.Add(numChannels)

	dbuser, ok := os.LookupEnv("BASEBALL_DB_USER")
	if ok == false {
		return errors.New("BASEBALL_DB_USER environment variable is not set")
	}
	dbpass, ok := os.LookupEnv("BASEBALL_DB_PASS")
	if ok == false {
		return errors.New("BASEBALL_DB_PASS environment variable is not set")
	}
	dbname, ok := os.LookupEnv("BASEBALL_DB_NAME")
	if ok == false {
		return errors.New("BASEBALL_DB_NAME environment variable is not set")
	}

	connectionString := fmt.Sprintf("user='%s' password='%s' dbname='%s' sslmode='%s'",
		dbuser, dbpass, dbname, "disable")
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("Unable to ping the database")
		return err
	}
	dbO.db = db

	dbO.tables = make(map[string]bool)

	return nil
}

/*
Stop the pipeline stage in a graceful manner
*/
func (dbO *DatabaseOutput) Stop() {
	for _, channel := range dbO.DataInput {
		close(channel)
	}
	dbO.wg.Wait()
	err := dbO.db.Close()
	if err != nil {
		fmt.Println("Error when closing the baseball database")
	}
}

/*
Abort the pipeline stage immediately
*/
func (dbO *DatabaseOutput) Abort() {
	for _, channel := range dbO.DataInput {
		close(channel)
	}
	err := dbO.db.Close()
	if err != nil {
		fmt.Println("Error when closing the baseball database")
	}
}

/*
Run the pipeline stage
*/
func (dbO *DatabaseOutput) Run() {
	for _, channel := range dbO.DataInput {
		go dbO.runChannelInput(channel)
	}
}

func (dbO *DatabaseOutput) runChannelInput(input chan string) {
	defer dbO.wg.Done()
	for inputData := range input {
		dbO.wg.Add(1)
		dbO.updateRecord(inputData)
		dbO.wg.Done()
	}
}

func (dbO *DatabaseOutput) updateRecord(record string) {
	// Grab the record type from the JSON-formatted string
	if strings.HasPrefix(record, "{\"RecordName\":") == true {
		endOfType := strings.Index(record[15:], "\"") + 15
		recordType := record[15:endOfType]

		_, tableCreated := dbO.tables[recordType]
		if tableCreated == false {
			dbO.tables[recordType] = false
		}

		switch recordType {
		case "VenueRecord":
			var vR records.VenueRecord
			err := json.Unmarshal([]byte(record), &vR)
			if err != nil {
				fmt.Println("Unable to unmarshal VenueRecord")
			}
			if tableCreated == false {
				vR.CreateTable(dbO.db)
				dbO.tables[recordType] = true
			}
			vR.UpdateRecord(dbO.db)
		case "LeagueRecord":
			var lR records.LeagueRecord
			err := json.Unmarshal([]byte(record), &lR)
			if err != nil {
				fmt.Println("Unable to unmarshal League Record")
			}
			if tableCreated == false {
				lR.CreateTable(dbO.db)
				dbO.tables[recordType] = true
			}
			lR.UpdateRecord(dbO.db)
		case "DivisionRecord":
			var dR records.DivisionRecord
			err := json.Unmarshal([]byte(record), &dR)
			if err != nil {
				fmt.Println("Unable to unmarshal DivisionRecord")
			}
			if tableCreated == false {
				dR.CreateTable(dbO.db)
				dbO.tables[recordType] = true
			}
			dR.UpdateRecord(dbO.db)
		case "TeamRecord":
			var tR records.TeamRecord
			err := json.Unmarshal([]byte(record), &tR)
			if err != nil {
				fmt.Println("Unable to unmarshal TeamRecord")
			}
			if tableCreated == false {
				tR.CreateTable(dbO.db)
				dbO.tables[recordType] = true
			}
			tR.UpdateRecord(dbO.db)
		case "StandingRecord":
			var sR records.StandingRecord
			err := json.Unmarshal([]byte(record), &sR)
			if err != nil {
				fmt.Println("Unable to unmarshal StandingRecord")
			}
			if tableCreated == false {
				sR.CreateTable(dbO.db)
				dbO.tables[recordType] = true
			}
			sR.UpdateRecord(dbO.db)
		case "GameRecord":
			var gR records.GameRecord
			err := json.Unmarshal([]byte(record), &gR)
			if err != nil {
				fmt.Println("Unable to unmarshal GameRecord")
			}
			if tableCreated == false {
				gR.CreateTable(dbO.db)
				dbO.tables[recordType] = true
			}
			gR.UpdateRecord(dbO.db)
		case "GameStatusRecord":
			var gsR records.GameStatusRecord
			err := json.Unmarshal([]byte(record), &gsR)
			if err != nil {
				fmt.Println("Unable to unmarshal GameStatusRecord")
			}
			if tableCreated == false {
				gsR.CreateTable(dbO.db)
				dbO.tables[recordType] = true
			}
			gsR.UpdateRecord(dbO.db)
		case "InningScoreRecord":
			var isR records.InningScoreRecord
			err := json.Unmarshal([]byte(record), &isR)
			if err != nil {
				fmt.Println("Unable to unmarshal InningScoreRecord")
			}
			if tableCreated == false {
				isR.CreateTable(dbO.db)
				dbO.tables[recordType] = true
			}
			isR.UpdateRecord(dbO.db)
		default:
			fmt.Printf("Unexpected record type %s", recordType)
		}
	}
}
