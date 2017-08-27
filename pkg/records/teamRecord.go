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
TeamRecord is the specifc data record for each team in a league
*/
type TeamRecord struct {
	RecordName    string
	EffectiveDate time.Time
	ID            int64
	Name          string
	Code          string
	City          string
	LeagueID      int64
	Division      string
}

/*
ScreenOutput displays the record on the screen
*/
func (tR *TeamRecord) ScreenOutput() {
	fmt.Println(tR)
}

/*
FileOutput displays the record on the screen
*/
func (tR *TeamRecord) FileOutput(filePtr *os.File) {
	fmt.Fprintf(filePtr, "%s|%d|%s|%s|%s|%d|%s\n",
		tR.EffectiveDate.Format(time.UnixDate),
		tR.ID,
		tR.Name,
		tR.Code,
		tR.City,
		tR.LeagueID,
		tR.Division,
	)
}
