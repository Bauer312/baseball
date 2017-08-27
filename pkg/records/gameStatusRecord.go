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
	"fmt"
	"os"
	"time"
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
