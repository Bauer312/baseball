/*
	Copyright 2017 Brian Bauer

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, fOftware
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package pipelineStage

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/bauer312/baseball/pkg/records"
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
FileOutput contains the elements of a pipeline stage that will accept
	strings of data and print them to the screen
*/
type FileOutput struct {
	DataInput []chan string
	wg        sync.WaitGroup
	basePath  string
	files     map[string]*os.File
}

/*
Init the pipeline stage,
*/
func (fO *FileOutput) Init() error {
	numChannels := len(fO.DataInput)
	fO.wg.Add(numChannels)

	switch runtime.GOOS {
	case "darwin":
		usr, err := user.Current()
		if err != nil {
			fmt.Println("Unable to determine user storage location")
			return err
		}
		fO.basePath = filepath.Join(usr.HomeDir, "Library/Application Support/com.13fpl.baseball/")
	default:
		fO.basePath = "/var/lib/com.13fpl.baseball/"
	}

	err := os.MkdirAll(fO.basePath, os.ModePerm)
	if err != nil {
		fmt.Println("Unable to validate storage location")
		return err
	}

	fO.files = make(map[string]*os.File)

	return nil
}

/*
Stop the pipeline stage in a graceful manner
*/
func (fO *FileOutput) Stop() {
	for _, channel := range fO.DataInput {
		close(channel)
	}
	fO.wg.Wait()

	for _, filePtr := range fO.files {
		filePtr.Close()
	}
}

/*
Abort the pipeline stage immediately
*/
func (fO *FileOutput) Abort() {
	for _, channel := range fO.DataInput {
		close(channel)
	}
	for _, filePtr := range fO.files {
		filePtr.Close()
	}
}

/*
Run the pipeline stage
*/
func (fO *FileOutput) Run() {
	for _, channel := range fO.DataInput {
		go fO.runChannelInput(channel)
	}
}

func (fO *FileOutput) runChannelInput(input chan string) {
	defer fO.wg.Done()
	for inputData := range input {
		fO.wg.Add(1)
		fO.writeRecord(inputData)
		fO.wg.Done()
	}
}

func (fO *FileOutput) writeRecord(record string) {
	// Grab the record type from the JSON-formatted string
	if strings.HasPrefix(record, "{\"RecordName\":") == true {
		endOfType := strings.Index(record[15:], "\"") + 15
		recordType := record[15:endOfType]

		_, ok := fO.files[recordType]
		if ok == false {
			fO.openFile(recordType)
		}

		switch recordType {
		case "VenueRecord":
			var vR records.VenueRecord
			err := json.Unmarshal([]byte(record), &vR)
			if err != nil {
				fmt.Println("Unable to unmarshal VenueRecord")
			}
			vR.FileOutput(fO.files[recordType])
		case "LeagueRecord":
			var lR records.LeagueRecord
			err := json.Unmarshal([]byte(record), &lR)
			if err != nil {
				fmt.Println("Unable to unmarshal League Record")
			}
			lR.FileOutput(fO.files[recordType])
		case "DivisionRecord":
			var dR records.DivisionRecord
			err := json.Unmarshal([]byte(record), &dR)
			if err != nil {
				fmt.Println("Unable to unmarshal DivisionRecord")
			}
			dR.FileOutput(fO.files[recordType])
		case "TeamRecord":
			var tR records.TeamRecord
			err := json.Unmarshal([]byte(record), &tR)
			if err != nil {
				fmt.Println("Unable to unmarshal TeamRecord")
			}
			tR.FileOutput(fO.files[recordType])
		case "StandingRecord":
			var sR records.StandingRecord
			err := json.Unmarshal([]byte(record), &sR)
			if err != nil {
				fmt.Println("Unable to unmarshal StandingRecord")
			}
			sR.FileOutput(fO.files[recordType])
		case "GameRecord":
			var gR records.GameRecord
			err := json.Unmarshal([]byte(record), &gR)
			if err != nil {
				fmt.Println("Unable to unmarshal GameRecord")
			}
			gR.FileOutput(fO.files[recordType])
		case "GameStatusRecord":
			var gsR records.GameStatusRecord
			err := json.Unmarshal([]byte(record), &gsR)
			if err != nil {
				fmt.Println("Unable to unmarshal GameStatusRecord")
			}
			gsR.FileOutput(fO.files[recordType])
		case "InningScoreRecord":
			var isR records.InningScoreRecord
			err := json.Unmarshal([]byte(record), &isR)
			if err != nil {
				fmt.Println("Unable to unmarshal InningScoreRecord")
			}
			isR.FileOutput(fO.files[recordType])
		default:
			fmt.Printf("Unexpected record type %s", recordType)
		}
	}
}

func (fO *FileOutput) openFile(recordType string) {
	fileName := recordType + ".dat"
	filePath := path.Join(fO.basePath, fileName)
	ptr, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
	fO.files[recordType] = ptr
}
