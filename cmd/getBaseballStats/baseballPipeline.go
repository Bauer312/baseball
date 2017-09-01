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
	fO         *pipelineStage.FileOutput
	sO         *pipelineStage.ScreenOutput
	dbO        *pipelineStage.DatabaseOutput
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

	//Convert a date range to a set of paths
	bp.dP = &pipelineStage.DateToPath{
		DataInput: make(chan pipelineStage.DateInputParameters),
		BaseURL:   "http://gd2.mlb.com",
	}
	bp.dP.Init()
	go bp.dP.Run()

	//Load the scoreboard file represented by a path
	bp.sB = &pipelineStage.ScoreBoardFile{
		DataInput: bp.dP.DataOutput,
		BaseURL:   "http://gd2.mlb.com",
		Client:    &client,
	}
	bp.sB.Init()
	go bp.sB.Run()

	//Deal with the output based upon the requested output handling method
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
		bp.dbO = &pipelineStage.DatabaseOutput{
			DataInput: []chan string{
				bp.sB.GameFileOutout,
				bp.sB.DataOutput,
			},
		}
		bp.dbO.Init()
		go bp.dbO.Run()
	default:
		return errors.New("Unrecognized output type: " + bp.outputType)
	}

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
		bp.dbO.Stop()
	default:
		return errors.New("Unrecognized output type: " + bp.outputType)
	}

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