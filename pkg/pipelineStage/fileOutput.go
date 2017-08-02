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
	"os"
	"sync"
	"time"
)

/*
FileOutputParameters represents the data that comes into this pipeline stage
*/
type FileOutputParameters struct {
	FileName   string
	RecordDate time.Time
	DataRecord string
}

/*
FileName represents the output to be used for a specific file name
*/
type FileName struct {
	FileName string
	FilePtr  *os.File
}

/*
FileOutput contains the elements of the stage
*/
type FileOutput struct {
	DataInput chan FileOutputParameters
	files     []FileName
	wg        sync.WaitGroup
	rwg       sync.WaitGroup
}

/*
ChannelListener should be run in a goroutine and will receive data records
	on the input channel.  Each data record will be saved to the file
	specified.
*/
func (fO *FileOutput) ChannelListener() {
	for inputData := range fO.DataInput {
		fO.rwg.Add(1)
		for _, file := range fO.files {
			if file.FileName[:len(file.FileName)-4] == inputData.FileName {
				_, err := file.FilePtr.WriteString(inputData.DataRecord)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
		fO.rwg.Done()
	}
	fO.rwg.Wait()

	//Tell the pipeline we are done
	fO.wg.Done()
}

/*
Init will create all channels and other initialization needs.
	The DataInput channel is the output of any previous
	pipeline stage so it shouldn't be created here
*/
func (fO *FileOutput) Init(files []FileName) error {
	fO.wg.Add(1)
	fO.files = files
	return nil
}

/*
Stop will close the input channel, causing the Channel Listener to stop
*/
func (fO *FileOutput) Stop() {
	close(fO.DataInput)
	fO.wg.Wait()
	for _, file := range fO.files {
		file.FilePtr.Close()
	}
}
