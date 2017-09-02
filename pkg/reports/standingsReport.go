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
	"time"
	//Make sure we can use Postgres
	_ "github.com/lib/pq"
)

/*
StandingsReport represents the data necessary to create a standings report
	for a single team
*/
type StandingsReport struct {
	League     string
	Division   string
	Name       string
	Wins       int
	Losses     int
	WinningPct float32
}

/*
GetStandingsReport retrieves data as of the provided date, using the provided
	database connection.
*/
func GetStandingsReport(db *sql.DB, asOf time.Time) {
	statement := `SELECT lr.name "League", dr.name "Division", tr.name "Name", sr.wins "Wins", sr.losses "Losses"
	FROM StandingRecord sr
	JOIN TeamRecord tr ON
	tr.id = sr.teamid
	JOIN LeagueRecord lr ON
	lr.id = tr.leagueid
	JOIN DivisionRecord dr ON
	tr.division = dr.code
	WHERE (sr.effectivedate,sr.teamid) in
	(SELECT MAX(effectivedate), teamid FROM StandingRecord WHERE effectivedate <= $1 GROUP BY teamid)
	ORDER BY lr.name, tr.division, sr.wins desc;`

	rows, err := db.Query(statement, asOf)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	var currentLeague string
	var currentDivision string
	for rows.Next() {
		var newRecord StandingsReport
		if err := rows.Scan(&newRecord.League, &newRecord.Division, &newRecord.Name, &newRecord.Wins, &newRecord.Losses); err != nil {
			log.Fatal(err)
		}
		totalGames := newRecord.Wins + newRecord.Losses
		newRecord.WinningPct = float32(newRecord.Wins) / float32(totalGames)

		if currentLeague != newRecord.League {
			currentLeague = newRecord.League
		}
		if currentDivision != newRecord.Division {
			currentDivision = newRecord.Division
			fmt.Printf("\n%15s %s\n%10s %4s %4s %5s\n", currentLeague, currentDivision, "", "W", "L", "Pct")

		}
		fmt.Printf("%10s %4d %4d  %1.3f\n", newRecord.Name, newRecord.Wins, newRecord.Losses, newRecord.WinningPct)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
