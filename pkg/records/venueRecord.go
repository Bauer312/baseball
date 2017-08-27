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
VenueRecord is the specific data record for each venue
*/
type VenueRecord struct {
	RecordName    string
	EffectiveDate time.Time
	ID            int64
	Name          string
	Location      string
	Channel       string
}

/*
ScreenOutput displays the record on the screen
*/
func (vR *VenueRecord) ScreenOutput() {
	fmt.Println(vR)
}

/*
FileOutput displays the record on the screen
*/
func (vR *VenueRecord) FileOutput(filePtr *os.File) {
	fmt.Fprintf(filePtr, "%s|%d|%s|%s|%s\n",
		vR.EffectiveDate.Format(time.UnixDate),
		vR.ID,
		vR.Name,
		vR.Location,
		vR.Channel,
	)
}
