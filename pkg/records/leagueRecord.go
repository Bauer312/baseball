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
LeagueRecord is the specific data record for each league
*/
type LeagueRecord struct {
	RecordName    string
	EffectiveDate time.Time
	ID            int64
	Name          string
	SportCode     string
}

/*
ScreenOutput displays the record on the screen
*/
func (lR *LeagueRecord) ScreenOutput() {
	fmt.Println(lR)
}

/*
FileOutput displays the record on the screen
*/
func (lR *LeagueRecord) FileOutput(filePtr *os.File) {
	fmt.Fprintf(filePtr, "%s|%d|%s|%s\n",
		lR.EffectiveDate.Format(time.UnixDate),
		lR.ID,
		lR.Name,
		lR.SportCode,
	)
}
