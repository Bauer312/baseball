/*
	Copyright 2019 Brian Bauer

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

package datepath

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

/*
DatePath turns a single date path into the paths of all
	relevant files for that date
*/
type DatePath struct {
	DatePath chan string
	FilePath chan string
	reqWG    sync.WaitGroup
}

/*
ChannelListener should be run in a goroutine and will receive paths on the DatePath
	channel.  It will retrieve the data for that path and parse it, looking for all
	games in the file it receives.  For each of those games, it will publish the path
	to all the files that we care about.  To end this goroutine, close the DatePath
	channel.  When that happens, we wait for all outstanding URL requests to finish
	and the paths to be published on the FilePath channel.  Once that finishes, the
	goroutine exits.
*/
func (fP *DatePath) ChannelListener(client *http.Client) {
	for inputPath := range fP.DatePath {
		fP.reqWG.Add(1)
		fmt.Printf("\t\tRequesting %s\n", inputPath)
		resp, err := client.Get(inputPath)
		if err != nil {
			fmt.Println(err.Error())
		}
		fP.tokenize(inputPath, resp)
	}
	fmt.Println("FilePath input channel has closed; waiting for all requests to finish")
	fP.reqWG.Wait()
	fmt.Println("All requests have finished")
	close(fP.FilePath)
}

func (fP *DatePath) tokenize(dataPath string, resp *http.Response) {
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
			fP.reqWG.Done()
			return
		case token == html.StartTagToken:
			t := tokenizer.Token()

			isAnchor := t.Data == "a"
			if isAnchor {
				for _, a := range t.Attr {
					if a.Key == "href" {
						if i := strings.Index(a.Val, "gid_"); i >= 0 {
							gidPath := dataPath + a.Val[i:]
							if strings.HasSuffix(gidPath, "/") == false {
								gidPath = gidPath + "/"
							}

							fP.FilePath <- gidPath + "bis_boxscore.xml"
							fP.FilePath <- gidPath + "game.xml"
							fP.FilePath <- gidPath + "game_events.xml"
							fP.FilePath <- gidPath + "inning/inning_all.xml"
							fP.FilePath <- gidPath + "inning/inning_hit.xml"
						}
						break
					}
				}
			}
		}
	}
}

/*
Init will create all channels
*/
func (fP *DatePath) Init() {
	fP.DatePath = make(chan string)
	fP.FilePath = make(chan string)
}

/*
Done will close the DatePath channel, signalling that we are done
*/
func (fP *DatePath) Done() {
	close(fP.DatePath)
}
