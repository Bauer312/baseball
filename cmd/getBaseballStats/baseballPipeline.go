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
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Bauer312/baseball/pkg/pipelineStage"
)

/*
BaseballPipeline implements the pipeline interface and contains all of the
    pipeline stages needed to turn a date range into a set of flat files
    containing baseball stats
*/
type BaseballPipeline struct {
	dP         *pipelineStage.DateToPath
	sB         *pipelineStage.ScoreBoardFile
	dF         pipelineStage.DateFile
	gF         pipelineStage.GameFile
	fO         *pipelineStage.FileOutput
	sO         *pipelineStage.ScreenOutput
	outputType string
}

/*
Start will create and configure all of the stages of the pipeline
*/
func (bp *BaseballPipeline) Start(output string) error {
	fmt.Println("Starting the baseball pipeline")

	bp.outputType = output
	switch bp.outputType {
	case "screen":
		fmt.Println("Outputting all data to the screen")
	case "file":
		fmt.Println("Outputting all data to text files")
	case "db":
		fmt.Println("Outputting all data to the database")
	default:
		return errors.New("Unrecognized output type: " + bp.outputType)
	}

	//Reuse the same http client for all requests
	client := http.Client{Timeout: (10 * time.Second)}

	bp.dP = &pipelineStage.DateToPath{
		DataInput: make(chan pipelineStage.DateInputParameters),
		BaseURL:   "http://gd2.mlb.com",
	}
	bp.dP.Init()
	go bp.dP.Run()

	bp.sB = &pipelineStage.ScoreBoardFile{
		DataInput: bp.dP.DataOutput,
		BaseURL:   "http://gd2.mlb.com",
		Client:    &client,
	}
	bp.sB.Init()
	go bp.sB.Run()

	switch bp.outputType {
	case "screen":
		bp.sO = &pipelineStage.ScreenOutput{
			DataInput: []chan string{
				bp.sB.GameFileOutout,
				bp.sB.DataOutput,
			},
		}
		bp.sO.Init()
		go bp.sO.Run()
	case "file":
		bp.fO = &pipelineStage.FileOutput{
			DataInput: []chan string{
				bp.sB.GameFileOutout,
				bp.sB.DataOutput,
			},
		}
		bp.fO.Init()
		go bp.fO.Run()
	case "db":
	default:
		return errors.New("Unrecognized output type: " + bp.outputType)
	}

	bp.dF.Init()
	//bp.dF.DataInput = bp.dC.DataOutput

	bp.gF.Init()
	bp.gF.DataInput = bp.dF.GameFileOutout

	//bp.fO.Init(bp.FileNames(basePath))
	//bp.fO.Init()
	//bp.fO.DataInput = bp.gF.DataOutput

	// Start the pipelines in reverse order (why?)
	//go bp.fO.ChannelListener()
	go bp.gF.ChannelListener(&client)

	return nil
}

/*
End means that no more data will be sent into this pipeline
*/
func (bp *BaseballPipeline) End() error {
	bp.dP.Stop()
	bp.sB.Stop()

	switch bp.outputType {
	case "screen":
		bp.sO.Stop()
	case "file":
		bp.fO.Stop()
	case "db":
	default:
		return errors.New("Unrecognized output type: " + bp.outputType)
	}

	//Stop the gameFile stage
	bp.gF.Stop()

	fmt.Println("The baseball pipeline has shut down")

	return nil
}

/*
DateRange is the way in which we provide dates to the pipeline
*/
func (bp *BaseballPipeline) DateRange(beg, end string) error {
	fmt.Printf("Using a date range of %s - %s\n", beg, end)

	data := pipelineStage.DateInputParameters{Beg: beg, End: end}
	bp.dP.DataInput <- data

	return nil
}

/*
FileNames creates a bunch of file pointers that will be used
	to output data to files
*/
/*
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
*/
