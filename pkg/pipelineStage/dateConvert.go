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

package pipelineStage

import (
	"fmt"
	"time"

	"github.com/Bauer312/baseball/pkg/pipeline"
)

/*
DateConvert contains the elements of the stage
*/
type DateConvert struct {
	DataInput  chan time.Time
	DataOutput chan string
	Control    pipeline.StageControl
}

/*
ChannelListener should be run in a goroutine and will receive data on the input channel
	and the input control channel.  The time elements will be converted to URLs in the
	form of strings
*/
func (dC *DateConvert) ChannelListener(baseURL string) {
	var controlString string
	var inputDate time.Time
	for {
		select {
		case inputDate = <-dC.DataInput:
			year := inputDate.Year()
			month := inputDate.Month()
			day := inputDate.Day()
			dC.DataOutput <- fmt.Sprintf("%s/year_%04d/month_%02d/day_%02d/", baseURL, year, month, day)
		case controlString = <-dC.Control.Input:
			if controlString == "quit" {
				dC.Control.Output <- "done"
			}
		}
	}
}
