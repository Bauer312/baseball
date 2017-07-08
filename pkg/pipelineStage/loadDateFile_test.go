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

import "testing"

func TestLoadingDateFile(t *testing.T) {
	var testURL = []struct {
		InputData  URLLoadParameters
		OutputData string
	}{
		{
			InputData:  URLLoadParameters{URL: "http://www.test.com/a/b/c/year_2017/month_01/day_01/"},
			OutputData: "http://www.test.com/a/b/c/year_2017/month_01/day_01/",
		},
		{
			InputData:  URLLoadParameters{URL: "http://www.test.com/a/b/c/year_2017/month_01/day_02/"},
			OutputData: "http://www.test.com/a/b/c/year_2017/month_01/day_02/",
		},
	}

	for _, ex := range testURL {
		var uL URLLoad
		uL.Init()
		// DataInput channels don't get created automatically
		uL.DataInput = make(chan URLLoadParameters)

		go uL.ChannelListener()

		// Start the anonymous function that receives the output of the method under test
		go func(expected string) {
			output := <-uL.DataOutput
			if output != expected {
				t.Errorf("Output mismatch: Expected %s but received %s\n", output, expected)
			}

			// All data has been received, go ahead and send the signal that will cause the method under test to return
			uL.Control.Input <- "quit"
		}(ex.OutputData)

		// Send the input data to the input channel
		uL.DataInput <- ex.InputData

		// Wait until the method under test returns
		<-uL.Control.Output
	}
}
