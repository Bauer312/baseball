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
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestLoadingDateFile(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "<!DOCTYPE html><html><head></head><body>")
		fmt.Fprintln(w, "<a href=\"gid_test\">gid_test</a>")
		fmt.Fprintln(w, "</body></html>")
		return
	}))
	defer ts.Close()

	var testURL = []struct {
		InputData      string
		OutputData     []string
		GameFileOutput []string
	}{
		{
			InputData: ts.URL,
			OutputData: []string{
				ts.URL + "/gid_test/game_events.xml",
			},
			GameFileOutput: []string{
				ts.URL + "/gid_test/game.xml",
			},
		},
	}

	type ctrl struct {
		DF DateFile
		WG sync.WaitGroup
	}

	for _, ex := range testURL {
		data := ctrl{}

		data.DF.Init()
		// DataInput channels don't get created automatically
		data.DF.DataInput = make(chan string)
		go data.DF.ChannelListener(&http.Client{Timeout: (3 * time.Second)})

		// Start the anonymous function that receives the output of the method under test
		data.WG.Add(1)
		go func(expected []string) {
			for i, exOutput := range expected {
				output := <-data.DF.DataOutput
				if output != exOutput {
					t.Errorf("Output element %d mismatch: expected %s but received %s", i, exOutput, output)
				} else {
					t.Logf("Output element %d matched: %s == %s", i, exOutput, output)
				}
			}
			data.WG.Done()
		}(ex.OutputData)

		// Start the anonymous function that receives the GameFile output
		data.WG.Add(1)
		go func(expected []string) {
			for i, exOutput := range expected {
				output := <-data.DF.GameFileOutout
				if output != exOutput {
					t.Errorf("Output element %d mismatch: expected %s but received %s", i, exOutput, output)
				} else {
					t.Logf("Output element %d matched: %s == %s", i, exOutput, output)
				}
			}
			data.WG.Done()
		}(ex.GameFileOutput)

		// Send the input data to the input channel
		data.DF.DataInput <- ex.InputData

		data.DF.Stop()

		data.WG.Wait()
	}

}
