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
InningScoreRecord is the specific data about an inning in a game
*/
type InningScoreRecord struct {
	RecordName    string
	EffectiveDate time.Time
	GameID        int64
	Inning        int
	AwayTeamRuns  int
	HomeTeamRuns  int
}

/*
ScreenOutput displays the record on the screen
*/
func (isR *InningScoreRecord) ScreenOutput() {
	fmt.Println(isR)
}

/*
FileOutput displays the record on the screen
*/
func (isR *InningScoreRecord) FileOutput(filePtr *os.File) {
	fmt.Fprintf(filePtr, "%s|%d|%d|%d|%d\n",
		isR.EffectiveDate.Format(time.UnixDate),
		isR.GameID,
		isR.Inning,
		isR.AwayTeamRuns,
		isR.HomeTeamRuns,
	)
}