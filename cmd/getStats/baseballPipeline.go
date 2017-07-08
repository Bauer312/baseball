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
	uL pipelineStage.DateFile
}

/*
Start will create and configure all of the stages of the pipeline
*/
func (bp *BaseballPipeline) Start() error {
	fmt.Println("Starting the baseball pipeline")

	bp.dI.Init()
	bp.dI.DataInput = make(chan pipelineStage.DateInputParameters)

	bp.dC.Init()
	bp.dC.DataInput = bp.dI.DataOutput

	bp.uL.Init()
	bp.uL.DataInput = bp.dC.DataOutput

	// Start the pipelines in reverse order (why?)
	go bp.uL.ChannelListener(&http.Client{Timeout: (10 * time.Second)})
	go bp.dC.ChannelListener("http://gd2.mlb.com/components/game/mlb")
	go bp.dI.ChannelListener()

	// Listen for the final stage to send output
	go func() {
		for {
			output := <-bp.uL.DataOutput
			fmt.Println(output)
		}

	}()

	return nil
}

/*
End means that no more data will be sent into this pipeline
*/
func (bp *BaseballPipeline) End() error {
	//Stop the dateInput stage
	close(bp.dI.DataInput)
	<-bp.dI.Control.Output

	//Stop the dateConvert stage
	close(bp.dC.DataInput)
	<-bp.dC.Control.Output

	//Stop the urlLoad stage
	close(bp.uL.DataInput)
	<-bp.uL.Control.Output

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
