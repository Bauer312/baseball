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

package reports

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path"

	//Make sure we can use Postgres
	_ "github.com/lib/pq"
)

/*
StandingsReport represents the data necessary to create a standings report
	for a single team
*/
type StandingsReport struct {
	Name       string
	Wins       int
	Losses     int
	WinningPct float32
}

/*
GetStandingsReport retrieves data as of the provided date, using the provided
	database connection.
*/
func GetStandingsReport(db *sql.DB, asOf, league, division, output string) {
	statement := `SELECT tr.name "Name", sr.wins "Wins", sr.losses "Losses"
	FROM StandingRecord sr
	JOIN TeamRecord tr ON
	tr.id = sr.teamid
	JOIN LeagueRecord lr ON
	lr.id = tr.leagueid
	JOIN DivisionRecord dr ON
	tr.division = dr.code
	WHERE (sr.effectivedate, sr.teamid) in
	(SELECT MAX(effectivedate), teamid FROM StandingRecord WHERE effectivedate <= $1 GROUP BY teamid)
	AND lr.name = $2 AND dr.name = $3
	ORDER BY sr.wins desc;`

	rows, err := db.Query(statement, asOf, league, division)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	var outputString string
	priorRecord := false
	if len(output) == 0 {
		fmt.Printf("\n%15s %s\n%10s %4s %4s %5s\n", league, division, "", "W", "L", "Pct")
	} else {
		outputString = fmt.Sprintf("{\n\t\"league\":\"%s\",\n\t\"division\":\"%s\",\n\t\"records\":[\n", league, division)
	}
	for rows.Next() {
		var newRecord StandingsReport
		if err := rows.Scan(&newRecord.Name, &newRecord.Wins, &newRecord.Losses); err != nil {
			log.Fatal(err)
		}
		totalGames := newRecord.Wins + newRecord.Losses
		newRecord.WinningPct = float32(newRecord.Wins) / float32(totalGames)
		if len(output) == 0 {
			fmt.Printf("%10s %4d %4d  %1.3f\n", newRecord.Name, newRecord.Wins, newRecord.Losses, newRecord.WinningPct)
		} else {
			if priorRecord == true {
				outputString = fmt.Sprintf("%s,\n", outputString)
			}
			outputString = fmt.Sprintf("%s\t\t{\"team\":\"%s\",\"wins\":\"%d\",\"losses\":\"%d\",\"pct\":\"%1.3f\"}", outputString, newRecord.Name, newRecord.Wins, newRecord.Losses, newRecord.WinningPct)
		}
		priorRecord = true
	}
	if len(output) > 0 {
		outputString = fmt.Sprintf("%s\t]\n}\n", outputString)

		switch league {
		case "American League":
			output = path.Join(output, "american")
		case "National League":
			output = path.Join(output, "national")
		}
		fileName := division + ".json"
		filePath := path.Join(output, fileName)
		ptr, err := os.Create(filePath)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer ptr.Close()
		fmt.Fprintf(ptr, "%s", outputString)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
