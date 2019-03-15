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
	"sync"
)

/*
ScreenOutput contains the elements of a pipeline stage that will accept
	strings of data and print them to the screen
*/
type ScreenOutput struct {
	DataInput []chan string
	wg        sync.WaitGroup
}

/*
Init the pipeline stage,
*/
func (sO *ScreenOutput) Init() error {
	numChannels := len(sO.DataInput)
	sO.wg.Add(numChannels)
	return nil
}

/*
Stop the pipeline stage in a graceful manner
*/
func (sO *ScreenOutput) Stop() {
	for _, channel := range sO.DataInput {
		close(channel)
	}
	sO.wg.Wait()
}

/*
Abort the pipeline stage immediately
*/
func (sO *ScreenOutput) Abort() {
	for _, channel := range sO.DataInput {
		close(channel)
	}
}

/*
Run the pipeline stage
*/
func (sO *ScreenOutput) Run() {
	for _, channel := range sO.DataInput {
		go sO.runChannelInput(channel)
	}
}

func (sO *ScreenOutput) runChannelInput(input chan string) {
	defer sO.wg.Done()
	for inputData := range input {
		fmt.Println(inputData)
	}
}
