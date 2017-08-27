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
DivisionRecord is the specific data record for each division in a league
*/
type DivisionRecord struct {
	RecordName    string
	EffectiveDate time.Time
	Name          string
	Code          string
}

/*
ScreenOutput displays the record on the screen
*/
func (dR *DivisionRecord) ScreenOutput() {
	fmt.Println(dR)
}

/*
FileOutput displays the record on the screen
*/
func (dR *DivisionRecord) FileOutput(filePtr *os.File) {
	fmt.Fprintf(filePtr, "%s|%s|%s\n",
		dR.EffectiveDate.Format(time.UnixDate),
		dR.Name,
		dR.Code,
	)
}
