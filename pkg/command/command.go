package command

import "flag"

/*
Command is the interface that all commands must conform to
*/
type Command interface {
	SetFlags(*flag.FlagSet, map[string]*string)
	Execute(map[string]*string)
}
