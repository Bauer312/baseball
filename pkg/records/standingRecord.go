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
StandingRecord is the specific data record for a team standing in
	its division in its league
*/
type StandingRecord struct {
	RecordName        string
	EffectiveDate     time.Time
	TeamID            int64
	Wins              int
	Losses            int
	GamesBack         float32
	WildcardGamesBack float32
}

/*
ScreenOutput displays the record on the screen
*/
func (sR *StandingRecord) ScreenOutput() {
	fmt.Println(sR)
}

/*
FileOutput displays the record on the screen
*/
func (sR *StandingRecord) FileOutput(filePtr *os.File) {
	fmt.Fprintf(filePtr, "%s|%d|%d|%d|%f|%f\n",
		sR.EffectiveDate.Format(time.UnixDate),
		sR.TeamID,
		sR.Wins,
		sR.Losses,
		sR.GamesBack,
		sR.WildcardGamesBack,
	)
}
