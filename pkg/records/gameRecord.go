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
GameRecord is the specific data record for a game
*/
type GameRecord struct {
	RecordName       string
	EffectiveDate    time.Time
	ID               int64
	ResumeDate       string
	OriginalDate     string
	GameType         string
	Tiebreaker       string
	GameDay          string
	DoubleHeader     string
	GameNumber       int
	TBDFlag          string
	Interleague      string
	ScheduledInnings int
	Description      string
	VenueID          int64
	AwayTeamID       int64
	HomeTeamID       int64
}

/*
ScreenOutput displays the record on the screen
*/
func (gR *GameRecord) ScreenOutput() {
	fmt.Println(gR)
}

/*
FileOutput displays the record on the screen
*/
func (gR *GameRecord) FileOutput(filePtr *os.File) {
	fmt.Fprintf(filePtr, "%s|%d|%s|%s|%s|%s|%s|%s|%d|%s|%s|%d|%s|%d|%d|%d\n",
		gR.EffectiveDate.Format(time.UnixDate),
		gR.ID,
		gR.ResumeDate,
		gR.OriginalDate,
		gR.GameType,
		gR.Tiebreaker,
		gR.GameDay,
		gR.DoubleHeader,
		gR.GameNumber,
		gR.TBDFlag,
		gR.Interleague,
		gR.ScheduledInnings,
		gR.Description,
		gR.VenueID,
		gR.AwayTeamID,
		gR.HomeTeamID,
	)
}
