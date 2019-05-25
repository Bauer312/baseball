package command

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

/*
ExtractWeatherLink contains information to extract weather linking data
	from Baseball Savant data
*/
type ExtractWeatherLink struct {
	inputDir  string
	outputDir string
}

/*
SetFlags creates the flags that are needed for this functionality
*/
func (ewl *ExtractWeatherLink) SetFlags(fs *flag.FlagSet, cmdMap map[string]*string) {
	cmdMap["inputDir"] = fs.String("input", ".", "Directory containing Savant CSV files")
	cmdMap["outputDir"] = fs.String("output", ".", "Directory to write SQL data file")
}

/*
Execute runs the functionality that produces the data needed
*/
func (ewl *ExtractWeatherLink) Execute(cmdMap map[string]*string) {
	ewl.inputDir = *cmdMap["inputDir"]
	ewl.outputDir = *cmdMap["outputDir"]

	files, err := ioutil.ReadDir(ewl.inputDir)
	if err != nil {
		log.Fatal(err)
	}

	outputFile := filepath.Join(ewl.outputDir, "savant-import.sql")
	ofp, err := os.Create(outputFile)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if strings.HasSuffix(strings.ToLower(f.Name()), ".csv") {
			fmt.Println(f.Name())
			readSavantCSV(filepath.Join(ewl.inputDir, f.Name()), ofp)
		}
	}

}

func readSavantCSV(f string, o *os.File) {
	fp, err := os.Open(f)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()
	r := csv.NewReader(fp)

	header, err := r.Read()
	if err != nil {
		log.Fatal(err)
	}
	filter := make([]bool, len(header))
	for i, v := range header {
		switch v {
		case "game_date", "player_name", "game_type", "home_team",
			"away_team", "inning", "inning_topbot", "sv_id", "game_pk",
			"at_bat_number", "pitch_number":
			filter[i] = true
		default:
			filter[i] = false
		}
	}

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		firstElement := true
		for i, v := range record {
			if filter[i] == true {
				if firstElement == false {
					fmt.Fprintf(o, ",")
				} else {
					firstElement = false
				}
				fmt.Fprintf(o, "'%s'", v)
			}
		}
		fmt.Fprintln(o, "")
	}
}
