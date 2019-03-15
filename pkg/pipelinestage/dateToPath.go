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

package pipelinestage

import (
	"fmt"
	"sync"
	"time"

	"github.com/bauer312/baseball/pkg/dateslice"
)

/*
DateInputParameters represents the data that comes into this pipeline stage
*/
type DateInputParameters struct {
	Beg string
	End string
}

/*
DateToPath contains the elements of a pipeline stage that will accept a
	set of date inputDataeters and, for each date, a path to the page
    containing data for that date.  Input and output are byte arrays
    containing marshalled JSON data.
*/
type DateToPath struct {
	DataInput  chan DateInputParameters
	DataOutput chan string
	BaseURL    string
	wg         sync.WaitGroup
}

/*
Init the pipeline stage,
*/
func (dP *DateToPath) Init() error {
	dP.wg.Add(1)
	dP.DataOutput = make(chan string)
	return nil
}

/*
Stop the pipeline stage in a graceful manner
*/
func (dP *DateToPath) Stop() {
	close(dP.DataInput)
	dP.wg.Wait()
}

/*
Abort the pipeline stage immediately
*/
func (dP *DateToPath) Abort() {
	close(dP.DataInput)
}

/*
Run the pipeline stage
*/
func (dP *DateToPath) Run() {
	defer dP.wg.Done()
	// This loop runs until the input channel is closed
	for inputData := range dP.DataInput {
		var dates []time.Time
		if len(inputData.Beg) != 0 && len(inputData.End) != 0 {
			dates = dateslice.RangeString(inputData.Beg, inputData.End)
		} else if len(inputData.Beg) != 0 && len(inputData.End) == 0 {
			dates = dateslice.RangeString(inputData.Beg, inputData.Beg)
		} else {
			fmt.Println("Invalid Date Input")
			return
		}

		for _, date := range dates {
			year := date.Year()
			month := date.Month()
			day := date.Day()
			dP.DataOutput <- fmt.Sprintf("%s/components/game/mlb/year_%04d/month_%02d/day_%02d/", dP.BaseURL, year, month, day)
		}
	}
}
