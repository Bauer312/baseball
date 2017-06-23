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

	"github.com/Bauer312/baseball/pkg/dateslice"
	"github.com/Bauer312/baseball/pkg/pipeline"
)

/*
DateInputParameters represents the data that comes into this pipeline stage
*/
type DateInputParameters struct {
	Beg string
	End string
}

/*
DateInput contains the elements of the stage
*/
type DateInput struct {
	DataInput  chan DateInputParameters
	DataOutput chan time.Time
	Control    pipeline.StageControl
}

/*
ChannelListener should be run in a goroutine and will receive data on the input channel
	and the input control channel.  The parameters are converted to a slice of time
	elements and these elements are then sent out over the output channel.
*/
func (dI *DateInput) ChannelListener() {
	var inputData DateInputParameters
	var controlString string

	for {
		select {
		case inputData = <-dI.DataInput:
			if len(inputData.Beg) != 0 && len(inputData.End) != 0 {
				output := dateslice.RangeString(inputData.Beg, inputData.End)
				for _, od := range output {
					dI.DataOutput <- od
				}
			} else if len(inputData.Beg) != 0 && len(inputData.End) == 0 {
				output := dateslice.RangeString(inputData.Beg, inputData.Beg)
				for _, od := range output {
					dI.DataOutput <- od
					//dI.DataOutput <- od
				}
			} else {
				fmt.Println("Invalid Date Input")
			}
		case controlString = <-dI.Control.Input:
			if controlString == "quit" {
				dI.Control.Output <- "done"
			}
		}
	}
}
