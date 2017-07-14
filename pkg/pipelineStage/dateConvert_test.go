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
	"testing"
	"time"
)

func TestDateConversion(t *testing.T) {
	var conversionTest = []struct {
		InputData  []time.Time
		OutputData []string
		BaseURL    string
	}{
		{
			[]time.Time{
				time.Date(2017, time.January, 1, 5, 0, 0, 0, time.UTC),
				time.Date(2017, time.January, 2, 5, 0, 0, 0, time.UTC),
			},
			[]string{
				"http://test.com/a/b/c/year_2017/month_01/day_01/",
				"http://test.com/a/b/c/year_2017/month_01/day_02/",
			},
			"http://test.com/a/b/c",
		},
		{
			[]time.Time{
				time.Date(2017, time.May, 29, 5, 0, 0, 0, time.UTC),
				time.Date(2017, time.May, 30, 5, 0, 0, 0, time.UTC),
				time.Date(2017, time.May, 31, 5, 0, 0, 0, time.UTC),
				time.Date(2017, time.June, 1, 5, 0, 0, 0, time.UTC),
			},
			[]string{
				"http://test.com/a/b/c/year_2017/month_05/day_29/",
				"http://test.com/a/b/c/year_2017/month_05/day_30/",
				"http://test.com/a/b/c/year_2017/month_05/day_31/",
				"http://test.com/a/b/c/year_2017/month_06/day_01/",
			},
			"http://test.com/a/b/c",
		},
	}

	for caseNumber, ex := range conversionTest {
		var dC DateConvert
		dC.Init()
		// DataInput channels don't get created automatically
		dC.DataInput = make(chan time.Time)

		// Start the method under test
		go dC.ChannelListener(ex.BaseURL)

		// Start the anonymous function that receives the output of the method under test
		go func(expected []string) {
			for i, exOutput := range expected {
				output := <-dC.DataOutput
				if output != exOutput {
					t.Errorf("Output element %d mismatch: expected %s but received %s", i, exOutput, output)
				} else {
					t.Logf("Output element %d matched: %s == %s", i, exOutput, output)
				}
			}

			//Check to ensure that the output channel is empty
			select {
			case <-dC.DataOutput:
				t.Errorf("Test Case %d received too many elements: expected %d but received at least %d",
					caseNumber,
					len(expected),
					len(expected)+1)
			default:
				break
			}
		}(ex.OutputData)

		// Send the input data to the input channel
		for _, data := range ex.InputData {
			dC.DataInput <- data
		}

		// Terminate by closing the data input channel
		close(dC.DataInput)

		// Wait until the method under test returns
		<-dC.Control.Output
	}
}
