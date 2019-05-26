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
			command.MainGameday()
		case "weather":
			cmdStruct = &command.ExtractWeatherLink{}
		case "loadsavant":
			cmdStruct = &command.LoadSavantData{}
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
}
