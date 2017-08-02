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

package main

import (
	"fmt"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/Bauer312/baseball/pkg/pipelineStage"
)

/*
BaseballPipeline implements the pipeline interface and contains all of the
    pipeline stages needed to turn a date range into a set of flat files
    containing baseball stats
*/
type BaseballPipeline struct {
	dI pipelineStage.DateInput
	dC pipelineStage.DateConvert
	dF pipelineStage.DateFile
	gF pipelineStage.GameFile
	fO pipelineStage.FileOutput
	wg sync.WaitGroup
}

/*
Start will create and configure all of the stages of the pipeline
*/
func (bp *BaseballPipeline) Start() error {
	fmt.Println("Starting the baseball pipeline")

	var basePath string
	switch runtime.GOOS {
	case "darwin":
		usr, _ := user.Current()
		basePath = filepath.Join(usr.HomeDir, "Library/Application Support/com.13fpl.baseball/")
	default:
		basePath = "/var/lib/com.13fpl.baseball/"
	}

	bp.dI.Init()
	bp.dI.DataInput = make(chan pipelineStage.DateInputParameters)

	bp.dC.Init()
	bp.dC.DataInput = bp.dI.DataOutput

	bp.dF.Init()
	bp.dF.DataInput = bp.dC.DataOutput

	bp.gF.Init()
	bp.gF.DataInput = bp.dF.GameFileOutout

	bp.fO.Init(bp.FileNames(basePath))
	bp.fO.DataInput = bp.gF.DataOutput

	//Reuse the same http client for all requests
	client := http.Client{Timeout: (10 * time.Second)}
	// Start the pipelines in reverse order (why?)
	go bp.fO.ChannelListener()
	go bp.gF.ChannelListener(&client)
	go bp.dF.ChannelListener(&client)
	go bp.dC.ChannelListener("http://gd2.mlb.com/components/game/mlb")
	go bp.dI.ChannelListener()

	// Listen for the final stage to send output
	go bp.PrintData()

	return nil
}

/*
PrintData will print the final output to the screen
*/
func (bp *BaseballPipeline) PrintData() {
	for output := range bp.dF.DataOutput {
		fmt.Println(output)
	}
	bp.wg.Done()
}

/*
End means that no more data will be sent into this pipeline
*/
func (bp *BaseballPipeline) End() error {
	//Stop the dateInput stage
	bp.dI.Stop()

	//Stop the dateConvert stage
	bp.dC.Stop()

	//Stop the urlLoad stage
	bp.dF.Stop()

	//Stop the gameFile stage
	bp.gF.Stop()

	//Stop the fileOutput stage
	bp.fO.Stop()

	//Stop the PrintData function
	bp.wg.Add(1)
	close(bp.dF.DataOutput)
	bp.wg.Wait()

	fmt.Println("The baseball pipeline has shut down")

	return nil
}

/*
DateRange is the way in which we provide dates to the pipeline
*/
func (bp *BaseballPipeline) DateRange(beg, end string) error {
	fmt.Printf("Using a date range of %s - %s\n", beg, end)

	data := pipelineStage.DateInputParameters{Beg: beg, End: end}
	bp.dI.DataInput <- data

	return nil
}

/*
FileNames creates a bunch of file pointers that will be used
	to output data to files
*/
func (bp *BaseballPipeline) FileNames(root string) []pipelineStage.FileName {
	files := []string{
		"teamInfo.dat",
		"gameInfo.dat",
		"stadiumInfo.dat",
	}
	retVal := make([]pipelineStage.FileName, len(files))

	err := os.MkdirAll(root, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	for i, file := range files {
		newPath := filepath.Join(root, file)
		fmt.Println(newPath)
		ptr, err := os.OpenFile(newPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		retVal[i] = pipelineStage.FileName{
			FileName: file,
			FilePtr:  ptr,
		}
	}

	return retVal
}
