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
