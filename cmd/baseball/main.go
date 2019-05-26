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
	"os"
	"strings"

	"github.com/bauer312/baseball/pkg/command"
)

func main() {
	if len(os.Args) == 1 {
		printCommands()
	} else {
		cmdMap := make(map[string]*string)
		fs := flag.NewFlagSet("cmd", flag.ContinueOnError)
		var cmdStruct command.Command

		cmd := strings.ToLower(os.Args[1])
		switch cmd {
		case "savant":
			cmdStruct = &command.GetSavantGames{}
		case "gameday":
			cmdStruct = &command.GetGamedayGames{}
		case "weather":
			cmdStruct = &command.ExtractWeatherLink{}
		case "loadsavant":
			cmdStruct = &command.LoadSavantData{}
		case "loadgameday":
			cmdStruct = &command.LoadGamedayData{}
		default:
			printCommands()
			return
		}

		cmdStruct.SetFlags(fs, cmdMap)

		fs.Parse(os.Args[2:])

		for k, v := range cmdMap {
			fmt.Println(k, *v)
		}

		cmdStruct.Execute(cmdMap)
	}
}

func printCommands() {
	fmt.Println("Available Commands")
	fmt.Println("\tsavant")
	fmt.Println("\tgameday")
	fmt.Println("\tweather")
	fmt.Println("\tloadsavant")
	fmt.Println("\tloadgameday")
}
