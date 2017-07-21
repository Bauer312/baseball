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
	"sync"
	"testing"
	"time"
)

func TestDateInputBegEnd(t *testing.T) {
	var begEndTest = []struct {
		InputData  DateInputParameters
		OutputData []time.Time
	}{
		{
			DateInputParameters{"20170101", "20170102"},
			[]time.Time{
				time.Date(2017, time.January, 1, 5, 0, 0, 0, time.UTC),
				time.Date(2017, time.January, 2, 5, 0, 0, 0, time.UTC),
			},
		},
		{
			DateInputParameters{"20170529", "20170601"},
			[]time.Time{
				time.Date(2017, time.May, 29, 5, 0, 0, 0, time.UTC),
				time.Date(2017, time.May, 30, 5, 0, 0, 0, time.UTC),
				time.Date(2017, time.May, 31, 5, 0, 0, 0, time.UTC),
				time.Date(2017, time.June, 1, 5, 0, 0, 0, time.UTC),
			},
		},
		{
			DateInputParameters{"20160227", "20160302"},
			[]time.Time{
				time.Date(2016, time.February, 27, 5, 0, 0, 0, time.UTC),
				time.Date(2016, time.February, 28, 5, 0, 0, 0, time.UTC),
				time.Date(2016, time.February, 29, 5, 0, 0, 0, time.UTC),
				time.Date(2016, time.March, 1, 5, 0, 0, 0, time.UTC),
				time.Date(2016, time.March, 2, 5, 0, 0, 0, time.UTC),
			},
		},
		{
			DateInputParameters{"20170101", ""},
			[]time.Time{
				time.Date(2017, time.January, 1, 5, 0, 0, 0, time.UTC),
			},
		},
	}

	type ctrl struct {
		DI DateInput
		WG sync.WaitGroup
	}

	for caseNumber, ex := range begEndTest {
		data := ctrl{}

		data.DI.Init()
		// DataInput channels don't get created automatically
		data.DI.DataInput = make(chan DateInputParameters)

		// Start the method under test
		go data.DI.ChannelListener()

		// Start the anonymous function that receives the output of the method under test
		data.WG.Add(1)
		go func() {
			for i, exOutput := range ex.OutputData {
				output := <-data.DI.DataOutput
				if output.Year() != exOutput.Year() {
					t.Errorf("Output element %d year mismatch: expected %d but received %d", i, exOutput.Year(), output.Year())
				}
				if output.Month() != exOutput.Month() {
					t.Errorf("Output element %d month mismatch: expected %d but received %d", i, int(exOutput.Month()), int(output.Month()))
				}
				if output.Day() != exOutput.Day() {
					t.Errorf("Output element %d day mismatch: expected %d but received %d", i, exOutput.Day(), output.Day())
				}
			}

			//Check to ensure that the output channel is empty
			select {
			case <-data.DI.DataOutput:
				t.Errorf("Test Case %d received too many elements: expected %d but received at least %d",
					caseNumber,
					len(ex.OutputData),
					len(ex.OutputData)+1)
			default:
				break
			}
			data.WG.Done()
		}()

		// Send the input data to the input channel
		data.DI.DataInput <- ex.InputData

		data.DI.Stop()

		//Wait until the anon function returns
		data.WG.Wait()
	}
}
