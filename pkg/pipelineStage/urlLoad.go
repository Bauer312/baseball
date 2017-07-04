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
	"net/http"
	"strings"

	"golang.org/x/net/html"

	"github.com/Bauer312/baseball/pkg/pipeline"
)

/*
URLLoad contains the elements of the stage
*/
type URLLoad struct {
	DataInput  chan string
	DataOutput chan string
	Control    pipeline.StageControl
}

/*
ChannelListener should be run in a goroutine and will receive data on the input channel
	and the input control channel.  The parameters are converted to a slice of time
	elements and these elements are then sent out over the output channel.
*/
func (uL *URLLoad) ChannelListener(client *http.Client) {
	for inputData := range uL.DataInput {
		resp, err := client.Get(inputData)
		if err != nil {
			uL.Control.Output <- err.Error()
		}
		uL.tokenize(resp)
	}
	uL.Control.Output <- "ended"
}

func (uL *URLLoad) tokenize(resp *http.Response) {
	defer resp.Body.Close()
	tokenizer := html.NewTokenizer(resp.Body)
	for {
		token := tokenizer.Next()

		switch {
		case token == html.ErrorToken:
			//The end of the file
			return
		case token == html.StartTagToken:
			t := tokenizer.Token()

			isAnchor := t.Data == "a"
			if isAnchor {
				for _, a := range t.Attr {
					if a.Key == "href" {
						if strings.HasPrefix(a.Val, "gid_") {
							uL.DataOutput <- a.Val
						}
						break
					}
				}
			}
		}
	}
}

/*
Init will create all channels and other initialization needs.
	The DataInput channel is the output of any previous
	pipeline stage so it shouldn't be created here
*/
func (uL *URLLoad) Init() error {
	uL.Control.Input = make(chan string)
	uL.Control.Output = make(chan string)

	uL.DataOutput = make(chan string, 5)

	return nil
}
