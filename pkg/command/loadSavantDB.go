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

package command

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/bauer312/baseball/pkg/db"
)

/*
LoadSavantData contains information to save
	Baseball Savant data
*/
type LoadSavantData struct {
	inputDir string
}

/*
SetFlags creates the flags that are needed for this functionality
*/
func (ewl *LoadSavantData) SetFlags(fs *flag.FlagSet, cmdMap map[string]*string) {
	cmdMap["inputDir"] = fs.String("input", ".", "Directory containing Savant CSV files")
}

/*
Execute runs the functionality that produces the data needed
*/
func (ewl *LoadSavantData) Execute(cmdMap map[string]*string) {
	ewl.inputDir = *cmdMap["inputDir"]

	files, err := ioutil.ReadDir(ewl.inputDir)
	if err != nil {
		log.Fatal(err)
	}

	bbdb := db.BaseballDB{}
	err = bbdb.ConnectToPostgres("BBALL")
	if err != nil {
		log.Fatal(err)
	}
	defer bbdb.Close()

	err = bbdb.ConfirmSavantMaster()
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if strings.HasSuffix(strings.ToLower(f.Name()), ".csv") {
			fmt.Println(f.Name())
			err = bbdb.LoadSavantCSV(filepath.Join(ewl.inputDir, f.Name()))
			if err != nil {
				log.Println(err)
			}
		}
	}

}
