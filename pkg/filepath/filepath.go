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

package filepath

import (
	"fmt"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

/*
FilePath turns a single date path into the paths of all
	relevant files for that date
*/
type FilePath struct {
	FilePath chan string
	basePath string
	wg       sync.WaitGroup
}

/*
ChannelListener should be run in a goroutine and will receive paths on the FilePath
	channel.  It will retrieve the data for that path and save it to a file in the
	location specified.  Once the channel is closed, the goroutine exits
*/
func (fP *FilePath) ChannelListener(client *http.Client) {
	for inputPath := range fP.FilePath {
		fmt.Printf("\tRequesting %s\n", inputPath)
		fP.wg.Add(1)
		time.Sleep(2 * time.Second)
		resp, err := client.Get(inputPath)
		if err != nil {
			fmt.Println(err.Error())
		} else if i := strings.Index(inputPath, "gid_"); i >= 0 {
			outputPath := filepath.Join(fP.basePath, strings.Replace(inputPath[i:], "/", "_", -1))
			fP.writeFile(outputPath, resp)
		}
		fP.wg.Done()
	}
}

func (fP *FilePath) writeFile(filePath string, resp *http.Response) {
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	resp.Write(f)
}

/*
Init will create all channels
*/
func (fP *FilePath) Init(output string) {
	fP.FilePath = make(chan string)
	usr, err := user.Current()
	if err != nil {
		fmt.Println("Unable to determine user storage location")
	}
	fP.basePath = filepath.Join(usr.HomeDir, "baseball/gameday/raw/")
	err = os.MkdirAll(fP.basePath, 0740)
	if err != nil {
		fmt.Println("Unable to validate storage location")
	}
}

/*
Done will close the DatePath channel, signalling that we are done
*/
func (fP *FilePath) Done() {
	fP.wg.Wait()
	close(fP.FilePath)
}
