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

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func main() {
	output := flag.String("output", "", "Output location for parsed flat files")
	input := flag.String("input", "", "Input location for raw xml files")

	flag.Parse()

	var inputPath string
	var outputPath string
	var err error
	if len(*output) == 0 {
		usr, _ := user.Current()
		homeDir := usr.HomeDir
		outputPath = filepath.Join(homeDir, "baseball")
		outputPath = filepath.Join(outputPath, "gameday")
		outputPath = filepath.Join(outputPath, "parsed")
	} else {
		outputPath, err = filepath.Abs(*output)
		if err != nil {
			log.Fatal(err)
		}
	}
	if len(*input) == 0 {
		usr, _ := user.Current()
		homeDir := usr.HomeDir
		inputPath = filepath.Join(homeDir, "baseball")
		inputPath = filepath.Join(inputPath, "gameday")
		inputPath = filepath.Join(inputPath, "raw")
	} else {
		inputPath, err = filepath.Abs(*input)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf("Raw data will be read from %s\n", inputPath)
	fmt.Printf("Parsed data will be saved into %s\n", outputPath)

	files, err := ioutil.ReadDir(inputPath)
	if err != nil {
		log.Fatal(err)
	}

	err = os.MkdirAll(outputPath, 0740)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if strings.Contains(file.Name(), "_bis_boxscore.xml") {
			fmt.Printf("Boxscore: %s\n", file.Name())
		} else if strings.Contains(file.Name(), "_game.xml") {
			fmt.Printf("Game: %s\n", file.Name())
		} else if strings.Contains(file.Name(), "_game_events.xml") {
			fmt.Printf("Game Events: %s\n", file.Name())
		} else if strings.Contains(file.Name(), "_inning_inning_all.xml") {
			fmt.Printf("Inning All: %s\n", file.Name())
		} else if strings.Contains(file.Name(), "_inning_innint_hit.xml") {
			fmt.Printf("Inning Hit: %s\n", file.Name())
		}
		//fmt.Println(file.Name())
	}
}
