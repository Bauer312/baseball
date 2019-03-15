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
	"net/http"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

/*
DateFile contains the elements of the stage
*/
type DateFile struct {
	DataInput      chan string
	DataOutput     chan string
	GameFileOutout chan string
	wg             sync.WaitGroup
	rwg            sync.WaitGroup
}

/*
ChannelListener should be run in a goroutine and will receive data on the input channel
	and the input control channel.  The parameters are converted to a slice of time
	elements and these elements are then sent out over the output channel.
*/
func (dF *DateFile) ChannelListener(client *http.Client) {
	for inputData := range dF.DataInput {
		dF.rwg.Add(1)
		resp, err := client.Get(inputData)
		if err != nil {
			fmt.Println(err.Error())
		}
		dF.tokenize(inputData, resp)
	}
	dF.rwg.Wait()

	//Tell the pipeline we are done
	dF.wg.Done()
}

func (dF *DateFile) tokenize(dataPath string, resp *http.Response) {
	defer resp.Body.Close()

	if strings.HasSuffix(dataPath, "/") == false {
		dataPath = dataPath + "/"
	}

	tokenizer := html.NewTokenizer(resp.Body)
	for {
		token := tokenizer.Next()

		switch {
		case token == html.ErrorToken:
			//The end of the file
			dF.rwg.Done()
			return
		case token == html.StartTagToken:
			t := tokenizer.Token()

			isAnchor := t.Data == "a"
			if isAnchor {
				for _, a := range t.Attr {
					if a.Key == "href" {
						if strings.HasPrefix(a.Val, "gid_") {
							gidPath := dataPath + a.Val
							if strings.HasSuffix(gidPath, "/") == false {
								gidPath = gidPath + "/"
							}
							gamePath := gidPath + "game.xml"
							dF.GameFileOutout <- gamePath
							gameEventsPath := gidPath + "game_events.xml"
							dF.DataOutput <- gameEventsPath
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
func (dF *DateFile) Init() error {
	dF.wg.Add(1)
	dF.DataOutput = make(chan string, 5)
	dF.GameFileOutout = make(chan string)

	return nil
}

/*
Stop will close the input channel, causing the Channel Listener to stop
*/
func (dF *DateFile) Stop() {
	close(dF.DataInput)
	dF.wg.Wait()
}
